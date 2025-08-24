package pipeline

import (
	"github.com/KhoalaS/godel/pkg/types"
)

type Node struct {
	Id    string `json:"id,omitempty"`
	Type  string `json:"type"`
	Phase Phase  `json:"phase"`
	Name  string `json:"name"`

	Error  string         `json:"error,omitempty"`
	Inputs []NodeInput    `json:"inputs,omitempty"`
	Status NodeStatus     `json:"status,omitempty"`
	Config map[string]any `json:"config,omitempty"`

	Run NodeFunc `json:"-"`
}

type NodeInput struct {
	Id       string   `json:"id"`
	Type     string   `json:"type"`
	Label    string   `json:"label"`
	Required bool     `json:"required"`
	Options  []string `json:"options,omitempty"` // for enums
}

type NodeFunc func(job *types.DownloadJob, data map[string]any) (types.DownloadJob, error)

type NodeStatus string

const (
	StatusPending NodeStatus = "pending"
	StatusRunning NodeStatus = "running"
	StatusSuccess NodeStatus = "success"
	StatusFailed  NodeStatus = "failed"
)
