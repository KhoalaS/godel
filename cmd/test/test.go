package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/KhoalaS/godel"
	"github.com/KhoalaS/godel/pkg/nodes"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	comm := make(chan pipeline.PipelineMessage, 96)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup
	var pipelines = make(chan *pipeline.Pipeline, 12)
	client := http.Client{}

	for i := range 4 {
		wg.Add(1)
		go godel.PipelineWorker(ctx, &wg, i, pipelines, &client)
	}

	p := pipeline.Pipeline{
		Id: uuid.NewString(),
		Nodes: []pipeline.Node{
			{
				Id:    uuid.NewString(),
				Type:  "limiter",
				Phase: pipeline.PrePhase,
				Name:  "Limiter",
				Config: map[string]any{
					"limit": "1000",
				},
				NodeType: pipeline.ConnectorNode,
				Status:   pipeline.StatusPending,
				Run:      nodes.LimiterNodeFunc,
			},
		},
		Comm: comm,
		Job: types.DownloadJob{
			Url: "http://localhost:9095",
		},
	}

	go func() {
		for msg := range comm {
			log.Debug().Any("msg", msg).Send()
		}
	}()

	log.Debug().Msg("adding pipeline to channel")

	pipelines <- &p

	<-ctx.Done()

	log.Info().Msg("Waiting for workers to exit")
	wg.Wait()

}
