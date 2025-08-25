package main

import (
	"os"

	"github.com/KhoalaS/godel/pkg/nodes"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	p := pipeline.Pipeline{
		Nodes: []pipeline.Node{
			{
				Id:    uuid.NewString(),
				Type:  "limiter",
				Phase: pipeline.PrePhase,
				Name:  "Limiter",
				Config: map[string]any{
					"limit": 1000,
				},
				Status: pipeline.StatusPending,
				Run:    nodes.LimiterNodeFunc,
			},
		},
		Comm: make(chan pipeline.PipelineMessage, 24),
	}

	job := types.DownloadJob{
		Url: "http://localhost:9095",
	}

	// simulate running the pipeline concurrently
	go p.Run(job)

	for message := range p.Comm {
		log.Debug().Any("msg", message).Send()
	}
}
