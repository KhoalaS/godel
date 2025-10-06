package pipeline

import (
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

var cdRegex = regexp.MustCompile(`filename="(.+?)"`)

func Download(ctx context.Context, client *http.Client, job *types.DownloadJob, pipeline IPipeline, nodeId string) (IFile, error) {
	if job.IsParent {
		log.Debug().Str("id", job.Id).Msg("Added bulk download")
		job.Status.Store(types.DOWNLOADING)
		//BroadCastUpdate(job)
		return nil, nil
	}

	parsedUrl, err := url.Parse(job.Url)
	if err != nil {
		return nil, err
	}

	log.Debug().Str("url", job.Url).Send()

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return nil, err
	}

	for k, v := range job.Headers {
		request.Header.Add(k, v)
	}

	if job.Status.Load() == types.PAUSED {
		log.Info().Int("bytes", job.BytesDownloaded).Msg("Partial file size")

		info, err := os.Stat(job.Filename)
		if err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				job.Status.Store(types.IDLE)
				job.BytesDownloaded = 0
				log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("File is missing on disk, restarting download")
			} else {
				return nil, err
			}
		} else {
			if job.BytesDownloaded != int(info.Size()) {
				log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("Missmatch between size on disk and stored size")
				job.BytesDownloaded = int(info.Size())
			}
			request.Header.Add("Range", fmt.Sprintf("bytes=%d-", job.BytesDownloaded))
		}
	}

	currentState := job.Status.Load()

	log.Debug().Str("url", job.Url).Msg("Making request")

	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if currentState != types.PAUSED && response.StatusCode != http.StatusOK {
		response.Body.Close()
		return nil, fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	if currentState == types.PAUSED && response.StatusCode != http.StatusPartialContent {
		response.Body.Close()
		return nil, fmt.Errorf("expected 206 Partial Content, got %d", response.StatusCode)
	}

	log.Info().Str("id", job.Id).Msg("Request successful")

	if strings.TrimSpace(job.Filename) == "" {
		job.Filename = nameFromContentDisposition(response)
		if job.Filename == "" {
			job.Filename = FallbackFilename(parsedUrl)
			if job.Filename == "" {
				job.Filename = uuid.NewString()
			}
		}
	}

	var outfile *os.File
	outPath := filepath.Join(job.DestPath, job.Filename)
	destDir := filepath.Dir(outPath)
	err = os.MkdirAll(destDir, 0755)
	if err != nil {
		return nil, err
	}

	if currentState == types.PAUSED {
		outfile, err = os.OpenFile(outPath, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return nil, err
		}
		log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Opened file for appending")
	} else {
		outfile, err = os.Create(outPath)
		if err != nil {
			return nil, err
		}
		log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Created new file")
	}

	defer outfile.Close()
	defer response.Body.Close()

	contentLength := response.Header.Get("content-length")

	if contentLength == "" {
		log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("No content length")
		contentLength = "-1"
	}

	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return nil, err
	}

	if currentState == types.PAUSED {
		contentLengthInt += job.BytesDownloaded
	}

	if job.Size == 0 && contentLengthInt > 0 {
		job.Size = contentLengthInt
	}

	bytesRead := job.BytesDownloaded
	buf := make([]byte, utils.CHUNK_SIZE)
	lastBytesRead := job.BytesDownloaded

	ticker := time.NewTicker(time.Second)
	done := make(chan bool)
	lastTs := time.Now()

	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				elapsed := time.Since(lastTs).Seconds()

				deltaBytes := bytesRead - lastBytesRead
				speed := float64(deltaBytes) / elapsed
				job.Speed = speed

				if contentLengthInt == -1 {
					log.Info().Str("filename", job.Filename).Str("id", job.Id).Msgf("Speed: %.2f MB/s (eta: unknown)", speed/1024/1024)
				} else {
					remaining := contentLengthInt - bytesRead
					eta := float64(remaining) / float64(deltaBytes) * float64(elapsed)
					job.Eta = eta

					log.Info().Str("filename", job.Filename).Str("id", job.Id).Msgf("Speed: %.2f MB/s (eta: %.2f seconds)", speed/1024/1024, eta)
				}

				lastBytesRead = bytesRead
				lastTs = time.Now()
				BroadCastUpdate(pipeline, NewProgressMessage(pipeline.GetId(), nodeId, float64(bytesRead)/float64(contentLengthInt)))
			}
		}
	}()

	var reader io.Reader

	if job.Limit > 0 {
		limit := rate.NewLimiter(rate.Limit(job.Limit), 2*job.Limit)
		reader = &utils.RateLimitReader{
			Limiter: limit,
			Reader:  response.Body,
			Ctx:     ctx,
		}
	} else {
		reader = response.Body
	}

	job.Status.Store(types.DOWNLOADING)

	log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Start downloading")

	defer close(done)
	defer updateParentJob(job)

	for {
		select {
		case <-ctx.Done():
			log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("download canceled")
			job.Status.Store(types.PAUSED)
			return nil, ctx.Err()
		case <-job.CancelCh:
			if job.DeleteOnCancel {
				outfile.Close()
				log.Info().Str("filename", job.Filename).Msg("Removing file")
				os.Remove(job.Filename)
			}
			job.Status.Store(types.CANCELED)
			return nil, errors.New("download canceled")
		case <-job.PauseCh:
			job.BytesDownloaded = bytesRead
			job.Status.Store(types.PAUSED)
			return nil, errors.New("download paused")
		default:
			n, err := reader.Read(buf)
			if n > 0 {
				_, writeErr := outfile.Write(buf[:n])
				if writeErr != nil {
					job.Status.Store(types.ERROR)
					return nil, writeErr
				}
				bytesRead += n
				job.BytesDownloaded = bytesRead
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Done")
					job.Status.Store(types.DONE)
					return &FileWrapper{
						file:              outfile,
						destinationFolder: job.DestPath,
					}, nil
				} else {
					job.Status.Store(types.ERROR)
					return nil, err
				}
			}
		}
	}
}

func nameFromContentDisposition(r *http.Response) string {
	cd := r.Header.Get("content-disposition")
	if cd == "" {
		return ""
	}

	m := cdRegex.FindStringSubmatch(cd)
	if len(m) != 2 {
		return ""
	}

	return m[1]
}

func updateParentJob(job *types.DownloadJob) {

	if job.ParentId == "" {
		return
	}

	if parentJob, ok := registries.JobRegistry.Load(job.ParentId); ok {
		log.Debug().Str("parentId", job.ParentId).Msg("Updating parent")
		parentJob.BytesDownloaded = parentJob.BytesDownloaded + 1
		if parentJob.BytesDownloaded == parentJob.Size {
			parentJob.Status.Store(types.DONE)
		}
	}
}

func BroadCastUpdate(pipeline IPipeline, message PipelineMessage) {
	log.Info().Msg("broadcasting to client")
	pipeline.SendMessage(message)
}

func FallbackFilename(_url *url.URL) string {
	return path.Base(_url.Path)
}
