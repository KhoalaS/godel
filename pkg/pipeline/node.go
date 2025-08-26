package pipeline

import (
	"context"

	"github.com/KhoalaS/godel/pkg/types"
)

type Node struct {
	Id       string         `json:"id,omitempty"`
	Type     string         `json:"type"`
	Phase    Phase          `json:"phase"`
	Name     string         `json:"name"`
	NodeType NodeType       `json:"nodeType"`
	Error    string         `json:"error,omitempty"`
	Inputs   []NodeInput    `json:"inputs,omitempty"`
	Status   NodeStatus     `json:"status,omitempty"`
	Config   map[string]any `json:"config,omitempty"`

	Run NodeFunc `json:"-"`
}

type NodeType string

const (
	SourceNode    NodeType = "source"
	ConnectorNode NodeType = "connector"
	EndNode       NodeType = "end"
)

type NodeInput struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Label    string   `json:"label"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"` // for enums
}

type NodeFunc func(ctx context.Context, job types.DownloadJob, node Node, comm chan<- PipelineMessage) (types.DownloadJob, error)

type NodeStatus string

const (
	StatusPending NodeStatus = "pending"
	StatusRunning NodeStatus = "running"
	StatusSuccess NodeStatus = "success"
	StatusFailed  NodeStatus = "failed"
)
