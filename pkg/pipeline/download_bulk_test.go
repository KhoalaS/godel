package pipeline

import (
	"context"
	"net/http"
	"os"
	"sync"
	"testing"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func TestDownloadBulk(t *testing.T) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	job := types.NewDownloadJob()
	job.Urls = []string{
		"http://localhost:9999/files/1.txt",
		"http://localhost:9999/files/2.txt",
	}
	ctx := context.TODO()
	client := &http.Client{}
	jobs := make(chan *types.DownloadJob, 12)

	var wg sync.WaitGroup
	wg.Add(1)
	go downloadWorker(ctx, &wg, 0, jobs, client)

	DownloadBulk(ctx, client, job, jobs)

	close(jobs)
	wg.Wait()
	log.Debug().Any("jobs", registries.JobRegistry.All()).Send()
}

func downloadWorker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan *types.DownloadJob, client *http.Client) {
	log.Debug().Int("id", id).Msg("Worker online")

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int("id", id).Msg("Done signal sent to worker")
			return
		case job, ok := <-jobs:
			if !ok {
				log.Warn().Int("id", id).Msg("Unexpected jobs channel closure")
				return
			}

			log.Debug().Int("id", id).Msg("Downloading using worker")
			err := Download(ctx, client, job, "1", "1")
			if err != nil {
				log.Err(err).Str("status", string(job.Status.Load().(types.DownloadState))).Str("filename", job.Filename).Str("id", job.Id).Msg("error during download")
			}
		}
	}
}
