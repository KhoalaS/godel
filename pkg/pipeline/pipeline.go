package pipeline

import (
	"context"

	"github.com/google/uuid"
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
	defer close(p.Comm)

	ready := findStartNodes(p.Graph)
	done := map[string]bool{}

	for len(ready) > 0 {
		node := ready[0]
		ready = ready[1:]

		ApplyInputs(p.Graph, node)
		if err := node.Run(ctx, *node, p.Comm, p.Id, node.Id); err != nil {
			p.Comm <- PipelineMessage{
				PipelineId: p.Id,
				NodeId:     node.Id,
				NodeType:   node.Type,
				Type:       StatusMessage,
				Data: MessageData{
					Status: StatusFailed,
					Error:  err.Error(),
				},
			}
			return err
		}
		done[node.Id] = true

		p.Comm <- PipelineMessage{
			PipelineId: p.Id,
			NodeId:     node.Id,
			NodeType:   node.Type,
			Type:       StatusMessage,
			Data: MessageData{
				Status: StatusSuccess,
			},
		}

		for _, next := range p.Graph.Outgoing[node.Id] {
			if allDepsDone(next, done, p.Graph.Incoming) {
				ready = append(ready, next)
			}
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
