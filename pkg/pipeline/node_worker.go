package pipeline

import (
	"context"
	"sync"

	"github.com/rs/zerolog/log"
)

type NodeWorkerError struct {
	Error  error
	NodeId string
}

func NodeWorker(ctx context.Context, wg *sync.WaitGroup, id int, pipeline *Pipeline, nodes chan *Node, doneChan chan<- string, errChan chan<- NodeWorkerError) {
	log.Debug().Int("id", id).Msg("Node worker online")

	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Debug().Int("id", id).Msg("Done signal sent to node worker")
			return
		case node, ok := <-nodes:
			if !ok {
				log.Warn().Int("id", id).Msg("Unexpected node channel closure")
				return
			}

			log.Debug().Int("id", id).Msg("Executing node using worker")

			ApplyInputs(pipeline.Graph, node)
			err := node.Run(ctx, *node, pipeline.Comm, pipeline.Id, node.Id)
			if err != nil {
				log.Err(err).Send()
				errChan <- NodeWorkerError{
					Error:  err,
					NodeId: node.Id,
				}
			} else {
				doneChan <- node.Id
			}
		}
	}
}
