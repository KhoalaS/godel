package utils

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/KhoalaS/godel/pkg/types"
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

	response, err := client.Do(request)
	if err != nil {
		return err
	}

	if response.StatusCode != http.StatusOK {
		defer response.Body.Close()
		errorBody, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}

		errorMsg := string(errorBody)
		return fmt.Errorf("failed request with status code %d and body: %s", response.StatusCode, errorMsg)
	}

	if strings.TrimSpace(job.Filename) == "" {
		segments := strings.Split(parsedUrl.Path, "/")
		job.Filename = segments[len(segments)-1]
	}

	outfile, err := os.Create(job.Filename)
	if err != nil {
		return err
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

	bytesRead := 0
	buf := make([]byte, CHUNK_SIZE)
	lastBytesRead := 0

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

	job.Status = types.DOWNLOADING

	for {
		select {
		case <-ctx.Done():
			log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("download canceled")
			close(done)
			job.Status = types.PAUSED
			return ctx.Err()
		case <-job.CancelCh:
			close(done)
			job.Status = types.CANCELED
			return errors.New("download canceled")
		default:
			n, err := reader.Read(buf)
			if n > 0 {
				_, writeErr := outfile.Write(buf[:n])
				if writeErr != nil {
					job.Status = types.ERROR
					return writeErr
				}
				bytesRead += n

			}

			if err != nil {
				defer close(done)

				if errors.Is(err, io.EOF) {
					log.Info().Str("filename", job.Filename).Str("id", job.Id).Msg("Done")
					job.Status = types.DONE
					return nil
				} else {
					job.Status = types.ERROR
					return err
				}
			}
		}
	}
}
