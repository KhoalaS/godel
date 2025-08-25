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
	Nodes []Node `json:"nodes"`
}

func (p *Pipeline) Run(job types.DownloadJob) types.DownloadJob {
	_job := job
	for _, node := range p.Nodes {
		var err error
		_job, err = node.Run(_job, node)
		if err != nil {
			log.Warn().Err(err).Send()
		}
	}

	return _job
}
