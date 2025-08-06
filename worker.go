package godel

import (
	"context"
	"net/http"
	"sync"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/rs/zerolog/log"
)

func DownloadWorker(ctx context.Context, wg *sync.WaitGroup, id int, jobs <-chan *types.DownloadJob, client *http.Client) {
	log.Info().Int("id", id).Msg("Worker online")

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Info().Int("id", id).Msg("Done signal sent to worker")
			return
		case job, ok := <-jobs:
			if !ok {
				log.Warn().Int("id", id).Msg("Unexpected jobs channel closure")
				continue
			}

			log.Info().Int("id", id).Msg("Downloading using worker")
			err := utils.Download(ctx, client, job, nil)
			if err != nil {
				log.Err(err).Str("status", string(job.Status.Load().(types.DownloadState))).Str("filename", job.Filename).Str("id", job.Id).Msg("error during download")
			}
		}
	}
}
