package utils

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
	"strconv"
	"strings"
	"time"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/time/rate"
)

func Download(ctx context.Context, client *http.Client, job *types.DownloadJob, headers map[string]string) error {
	parsedUrl, err := url.Parse(job.Url)
	if err != nil {
		return err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, parsedUrl.String(), nil)
	if err != nil {
		return err
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
				return err
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

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if currentState != types.PAUSED && response.StatusCode != http.StatusOK {
		response.Body.Close()
		return fmt.Errorf("unexpected status code %d", response.StatusCode)
	}

	if currentState == types.PAUSED && response.StatusCode != http.StatusPartialContent {
		response.Body.Close()
		return fmt.Errorf("expected 206 Partial Content, got %d", response.StatusCode)
	}

	log.Info().Str("id", job.Id).Msg("Request successful")

	if strings.TrimSpace(job.Filename) == "" {
		job.Filename = FallbackFilename(parsedUrl)
		if job.Filename == "" {
			job.Filename = uuid.NewString()
		}
	}

	var outfile *os.File

	if currentState == types.PAUSED {
		outfile, err = os.OpenFile(job.Filename, os.O_APPEND|os.O_WRONLY, 0)
		if err != nil {
			return err
		}
		log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Opened file for appending")
	} else {
		outfile, err = os.Create(job.Filename)
		if err != nil {
			return err
		}
		log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Created new file")
	}

	defer outfile.Close()
	defer response.Body.Close()

	contentLength := response.Header.Get("content-length")

	if contentLength == "" {
		log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("No content length do io.Copy")
		_, err = io.Copy(outfile, response.Body)
		if err != nil {
			return err
		}
		return nil
	}

	contentLengthInt, err := strconv.Atoi(contentLength)
	if err != nil {
		return err
	}

	if currentState == types.PAUSED {
		contentLengthInt += job.BytesDownloaded
	}

	bytesRead := job.BytesDownloaded
	buf := make([]byte, CHUNK_SIZE)
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
				speed := float64(deltaBytes) / 1024 / 1024 / elapsed

				remaining := contentLengthInt - bytesRead
				eta := float64(remaining) / float64(deltaBytes) * float64(elapsed)

				log.Info().Str("filename", job.Filename).Str("id", job.Id).Msgf("Speed: %.2f MB/s (eta: %.2f seconds)", speed, eta)

				lastBytesRead = bytesRead
				lastTs = time.Now()
			}
		}
	}()

	var reader io.Reader

	if job.Limit > 0 {
		limit := rate.NewLimiter(rate.Limit(job.Limit), job.Limit)
		reader = &RateLimitReader{
			limiter: limit,
			reader:  response.Body,
			ctx:     ctx,
		}
	} else {
		reader = response.Body
	}

	job.Status.Store(types.DOWNLOADING)

	log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Start downloading")

	defer close(done)

	for {
		select {
		case <-ctx.Done():
			log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("download canceled")
			job.Status.Store(types.PAUSED)
			return ctx.Err()
		case <-job.CancelCh:
			if job.DeleteOnCancel {
				outfile.Close()
				log.Info().Str("filename", job.Filename).Msg("Removing file")
				os.Remove(job.Filename)
			}
			job.Status.Store(types.CANCELED)
			return errors.New("download canceled")
		case <-job.PauseCh:
			job.BytesDownloaded = bytesRead
			job.Status.Store(types.PAUSED)
			return errors.New("download paused")
		default:
			n, err := reader.Read(buf)
			if n > 0 {
				_, writeErr := outfile.Write(buf[:n])
				if writeErr != nil {
					job.Status.Store(types.ERROR)
					return writeErr
				}
				bytesRead += n
				job.BytesDownloaded = bytesRead
			}

			if err != nil {
				if errors.Is(err, io.EOF) {
					log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Done")
					job.Status.Store(types.DONE)
					return nil
				} else {
					job.Status.Store(types.ERROR)
					return err
				}
			}
		}
	}
}

func FallbackFilename(_url *url.URL) string {
	return path.Base(_url.Path)
}
