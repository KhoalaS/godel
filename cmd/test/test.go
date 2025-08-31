package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"sync"

	"github.com/KhoalaS/godel"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	comm := make(chan pipeline.PipelineMessage, 96)

	pipeline.NodeRegistry["int-input"] = pipeline.CreateIntInputNode()
	pipeline.NodeRegistry["download"] = pipeline.CreateDownloadNode()
	pipeline.NodeRegistry["downloader"] = pipeline.CreateDownloaderNode()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup
	var pipelines = make(chan *pipeline.Pipeline, 12)
	client := http.Client{}

	for i := range 4 {
		wg.Add(1)
		go godel.PipelineWorker(ctx, &wg, i, pipelines, &client)
	}

	intInputNode := pipeline.CreateIntInputNode()
	intInputNode.Id = "1"

	downloaderNode := pipeline.CreateDownloaderNode()
	downloaderNode.Id = "2"

	downloadNode := pipeline.CreateDownloadNode()
	downloadNode.Id = "3"

	// mock inputs
	if in, ok := intInputNode.Io["output"]; ok {
		in.Value = 1000000
	}
	if in, ok := downloadNode.Io["filename"]; ok {
		in.Value = "1.bin"
	}
	if in, ok := downloadNode.Io["output_dir"]; ok {
		in.Value = "./"
	}
	if in, ok := downloaderNode.Io["url"]; ok {
		in.Value = "http://localhost:9999/files/random.bin"
	}

	// use the int input for the limit

	graph := pipeline.NewGraph()
	graph.Edges = append(graph.Edges, pipeline.Edge{
		Source:       "1",
		Target:       "3",
		SourceHandle: "output",
		TargetHandle: "limit",
	})

	graph.Edges = append(graph.Edges, pipeline.Edge{
		Source:       "2",
		Target:       "3",
		SourceHandle: "job",
		TargetHandle: "job",
	})

	graph.Nodes[intInputNode.Id] = &intInputNode
	graph.Nodes[downloadNode.Id] = &downloadNode
	graph.Nodes[downloaderNode.Id] = &downloaderNode

	graph.Incoming[downloadNode.Id] = []*pipeline.Node{&intInputNode, &downloaderNode}
	graph.Outgoing[intInputNode.Id] = []*pipeline.Node{&downloadNode}
	graph.Outgoing[downloaderNode.Id] = []*pipeline.Node{&downloadNode}

	p := pipeline.Pipeline{
		Id:    uuid.NewString(),
		Graph: graph,
		Comm:  comm,
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
