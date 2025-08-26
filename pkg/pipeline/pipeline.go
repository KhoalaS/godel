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
	Job             types.DownloadJob    `json:"job"`
	Nodes           []Node               `json:"nodes"`
	Comm            chan PipelineMessage `json:"-"`
}

type PipelineMessage struct {
	PipelineId string      `json:"pipelineId"`
	NodeId     string      `json:"nodeId"`
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

	for _, node := range p.Nodes {
		p.Comm <- PipelineMessage{
			PipelineId: p.Id,
			NodeId:     node.Id,
			Type:       StatusMessage,
			Data: MessageData{
				Status: StatusRunning,
			},
		}

		var err error
		p.Job, err = node.Run(ctx, p.Job, node, p.Comm)
		if err != nil {
			log.Warn().Err(err).Send()
			node.Status = StatusFailed
			p.Comm <- PipelineMessage{
				PipelineId: p.Id,
				NodeId:     node.Id,
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
				Type:       StatusMessage,
				Data: MessageData{
					Status: StatusSuccess,
				},
			}
		}

	}

	return nil
}
