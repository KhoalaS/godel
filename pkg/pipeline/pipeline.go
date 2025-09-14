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
	Type  MessageType  `json:"type"`
	Data  any          `json:"data"`
	Level MessageLevel `json:"level"`
}

type MessageLevel string

const (
	ErrorMessageLevel = "error"
	WarnMessageLevel  = "warn"
	InfoMessageLevel  = "info"
)

func NewErrorMessage(pId string, nodeId string, err error) PipelineMessage {
	return PipelineMessage{
		Type:  NodeUpdateMessage,
		Level: ErrorMessageLevel,
		Data: NodeMessageData{
			PipelineId: pId,
			NodeId:     nodeId,
			Error:      err.Error(),
			Status:     StatusFailed,
		},
	}
}

func NewProgressMessage(pId string, nodeId string, progress float64) PipelineMessage {
	return PipelineMessage{
		Type:  NodeUpdateMessage,
		Level: InfoMessageLevel,
		Data: NodeMessageData{
			PipelineId: pId,
			NodeId:     nodeId,
			Status:     StatusRunning,
			Progress:   progress,
		},
	}
}

func NewStatusMessage(pId string, nodeId string, status NodeStatus) PipelineMessage {
	return PipelineMessage{
		Type:  NodeUpdateMessage,
		Level: InfoMessageLevel,
		Data: NodeMessageData{
			PipelineId: pId,
			NodeId:     nodeId,
			Status:     status,
		},
	}
}

func NewPipelineDoneMessage(pId string) PipelineMessage {
	return PipelineMessage{
		Type:  PipelineUpdateMessage,
		Level: InfoMessageLevel,
		Data: PipelineMessageData{
			PipelineId: pId,
			Status:     PipelineStatusDone,
		},
	}
}

func NewPipelineStartMessage(pId string) PipelineMessage {
	return PipelineMessage{
		Type:  PipelineUpdateMessage,
		Level: InfoMessageLevel,
		Data: PipelineMessageData{
			PipelineId: pId,
			Status:     PipelineStatusStarted,
		},
	}
}

func NewPipelineFailedMessage(pId string, err error) PipelineMessage {
	return PipelineMessage{
		Type:  PipelineUpdateMessage,
		Level: ErrorMessageLevel,
		Data: PipelineMessageData{
			PipelineId: pId,
			Error:      err.Error(),
			Status:     PipelineStatusFailed,
		},
	}
}

type NodeMessageData struct {
	PipelineId string     `json:"pipelineId"`
	NodeId     string     `json:"nodeId"`
	Error      string     `json:"error,omitempty"`
	Progress   float64    `json:"progress,omitempty"`
	Status     NodeStatus `json:"status"`
}

type PipelineMessageData struct {
	PipelineId string         `json:"pipelineId"`
	Error      string         `json:"error,omitempty"`
	Status     PipelineStatus `json:"status"`
}

type PipelineStatus string

const (
	PipelineStatusStarted PipelineStatus = "started"
	PipelineStatusFailed  PipelineStatus = "failed"
	PipelineStatusDone    PipelineStatus = "done"
)

type MessageType string

const (
	PipelineUpdateMessage MessageType = "pipelineUpdate"
	NodeUpdateMessage     MessageType = "nodeUpdate"
)

func (p *Pipeline) Run(ctx context.Context) error {
	BroadCastUpdate(NewPipelineStartMessage(p.Id))

	if HasCycle(p.Graph) {
		err := errors.New("graph has a cycle")
		BroadCastUpdate(NewPipelineFailedMessage(p.Id, err))
		return err
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

	BroadCastUpdate(NewPipelineDoneMessage(p.Id))

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
