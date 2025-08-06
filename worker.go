package godel

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/rs/zerolog/log"
)

func DownloadWorker(ctx context.Context, id int, jobs <-chan *types.DownloadJob, client *http.Client) {
	for {
		select {
		case <-ctx.Done():
			log.Info().Int("id", id).Msg("Done signal send to worker")
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			log.Info().Int("id", id).Msg("Downloading using worker")
			err := utils.Download(ctx, client, job, nil)
			if err != nil {
				log.Err(err).Msg("error during download")
			}

		}
	}
}
