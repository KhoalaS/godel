package pipeline

import (
	"context"
	"errors"
	"sync"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type Pipeline struct {
	Id              string               `json:"id"`
	FailOnNodeError bool                 `json:"failOnNodeError"`
	Graph           *Graph               `json:"nodes"`
	Comm            chan PipelineMessage `json:"-"`
}

func NewPipeline(g *Graph, comm chan PipelineMessage) *Pipeline {
	return &Pipeline{
		Id:              uuid.NewString(),
		FailOnNodeError: false,
		Graph:           g,
		Comm:            comm,
	}
}

type PipelineMessage struct {
	PipelineId string      `json:"pipelineId"`
	NodeId     string      `json:"nodeId"`
	NodeType   string      `json:"nodeType,omitempty"`
	Type       MessageType `json:"type"`
	Data       MessageData `json:"data"`
}

func NewErrorMessage(pId string, nodeId string, err error) PipelineMessage {
	return PipelineMessage{
		PipelineId: pId,
		NodeId:     nodeId,
		Type:       ErrorMessage,
		Data: MessageData{
			Error:  err.Error(),
			Status: StatusFailed,
		},
	}
}

func NewProgressMessage(pId string, nodeId string, progress float64) PipelineMessage {
	return PipelineMessage{
		PipelineId: pId,
		NodeId:     nodeId,
		Type:       ProgressMessage,
		Data: MessageData{
			Status:   StatusRunning,
			Progress: progress,
		},
	}
}

func NewStatusMessage(pId string, nodeId string, status NodeStatus) PipelineMessage {
	return PipelineMessage{
		PipelineId: pId,
		NodeId:     nodeId,
		Type:       StatusMessage,
		Data: MessageData{
			Status: status,
		},
	}
}

type MessageData struct {
	Error    string     `json:"error,omitempty"`
	Progress float64    `json:"progress,omitempty"`
	Status   NodeStatus `json:"status"`
}

type MessageType string

const (
	ErrorMessage    MessageType = "error"
	ProgressMessage MessageType = "progress"
	StatusMessage   MessageType = "status"
)

func (p *Pipeline) Run(ctx context.Context) error {
	if HasCycle(p.Graph) {
		return errors.New("graph has a cycle")
	}

	done := map[string]bool{}
	errChan := make(chan NodeWorkerError, 1)
	var wg sync.WaitGroup

	nodes := make(chan *Node, len(p.Graph.Nodes))
	doneChan := make(chan string, len(p.Graph.Nodes))

	// context to allow canelling the nodes in the pool
	_ctx, cancel := context.WithCancel(ctx)

	defer func() {
		cancel()

		close(nodes)

		log.Info().Msg("Waiting for node workers to shutdown")
		wg.Wait()
	}()

	for i := range 4 {
		wg.Add(1)
		go NodeWorker(_ctx, &wg, i, p, nodes, doneChan, errChan)
	}

	for _, node := range findStartNodes(p.Graph) {
		BroadCastUpdate(NewStatusMessage(p.Id, node.Id, StatusRunning))
		nodes <- node
	}

	for len(done) < len(p.Graph.Nodes) {
		select {
		case id := <-doneChan:
			done[id] = true

			BroadCastUpdate(NewStatusMessage(p.Id, id, StatusSuccess))
			for _, next := range p.Graph.Outgoing[id] {
				if allDepsDone(next, done, p.Graph.Incoming) {
					BroadCastUpdate(NewStatusMessage(p.Id, next.Id, StatusRunning))
					nodes <- next
				}
			}
		case err := <-errChan:
			log.Error().Err(err.Error).Send()
			BroadCastUpdate(NewErrorMessage(p.Id, err.NodeId, err.Error))
			return err.Error
		case <-ctx.Done():
			return ctx.Err()
		}
	}

	return nil
}

func findStartNodes(graph *Graph) []*Node {
	nodes := []*Node{}

	for _, node := range graph.Nodes {
		if inc, ok := graph.Incoming[node.Id]; !ok || len(inc) == 0 {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func allDepsDone(node *Node, done map[string]bool, incoming map[string][]*Node) bool {
	for _, dep := range incoming[node.Id] {
		if !done[dep.Id] {
			return false
		}
	}
	return true
}
