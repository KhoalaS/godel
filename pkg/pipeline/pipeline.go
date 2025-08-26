package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

type Phase string

const (
	PrePhase      Phase = "pre"
	DownloadPhase Phase = "download"
	AfterPhase    Phase = "after"
)

type Pipeline struct {
	Id              string               `json:"id"`
	FailOnNodeError bool                 `json:"failOnNodeError"`
	Nodes           []Node               `json:"nodes"`
	Comm            chan PipelineMessage `json:"-"`
}

type PipelineMessage struct {
	PipelineId string      `json:"pipelineId"`
	NodeId     string      `json:"nodeId"`
	NodeType   string      `json:"nodeType"`
	Type       MessageType `json:"type"`
	Data       MessageData `json:"data"`
}

type MessageData struct {
	Error    string     `json:"error"`
	Progress float64    `json:"progress"`
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

	job := *types.NewDownloadJob()

	for _, node := range p.Nodes {
		p.Comm <- PipelineMessage{
			PipelineId: p.Id,
			NodeId:     node.Id,
			NodeType:   node.Type,
			Type:       StatusMessage,
			Data: MessageData{
				Status: StatusRunning,
			},
		}

		var err error
		job, err = node.Run(ctx, job, node, p.Comm)
		if err != nil {
			log.Warn().Err(err).Send()
			node.Status = StatusFailed
			p.Comm <- PipelineMessage{
				PipelineId: p.Id,
				NodeId:     node.Id,
				NodeType:   node.Type,
				Type:       ErrorMessage,
				Data: MessageData{
					Error:  err.Error(),
					Status: StatusFailed,
				},
			}
			if p.FailOnNodeError {
				return err
			}
		} else {
			p.Comm <- PipelineMessage{
				PipelineId: p.Id,
				NodeId:     node.Id,
				NodeType:   node.Type,
				Type:       StatusMessage,
				Data: MessageData{
					Status: StatusSuccess,
				},
			}
		}

	}

	return nil
}
