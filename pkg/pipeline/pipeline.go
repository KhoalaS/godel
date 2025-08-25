package pipeline

import (
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
	Nodes []Node               `json:"nodes"`
	Comm  chan PipelineMessage `json:"-"`
}

type PipelineMessage struct {
	NodeId string      `json:"id"`
	Type   MessageType `json:"type"`
	Data   MessageData `json:"data"`
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

func (p *Pipeline) Run(job types.DownloadJob) types.DownloadJob {
	_job := job
	for _, node := range p.Nodes {
		p.Comm <- PipelineMessage{
			NodeId: node.Id,
			Type:   StatusMessage,
			Data: MessageData{
				Status: StatusRunning,
			},
		}

		var err error
		_job, err = node.Run(_job, node, p.Comm)
		if err != nil {
			log.Warn().Err(err).Send()
			p.Comm <- PipelineMessage{
				NodeId: node.Id,
				Type:   ErrorMessage,
				Data: MessageData{
					Error:  err.Error(),
					Status: StatusFailed,
				},
			}
			break
		}

		p.Comm <- PipelineMessage{
			NodeId: node.Id,
			Type:   StatusMessage,
			Data: MessageData{
				Status: StatusSuccess,
			},
		}
	}

	close(p.Comm)
	return _job
}
