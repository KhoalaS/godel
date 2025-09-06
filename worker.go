package godel

import (
	"context"
	"net/http"
	"sync"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/rs/zerolog/log"
)

func PipelineWorker(ctx context.Context, wg *sync.WaitGroup, id int, pipelines chan *pipeline.Pipeline, client *http.Client) {
	log.Debug().Int("id", id).Msg("Worker online")

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int("id", id).Msg("Done signal sent to worker")
			return
		case pipeline, ok := <-pipelines:
			if !ok {
				log.Warn().Int("id", id).Msg("Unexpected jobs channel closure")
				return
			}

			log.Debug().Int("id", id).Msg("Downloading using worker")

			err := pipeline.Run(ctx)
			if err != nil {
				log.Err(err).Send()
			}
		}
	}
}
