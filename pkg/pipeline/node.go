package pipeline

import (
	"context"
)

type Node struct {
	Id       string             `json:"id,omitempty"`
	Type     string             `json:"type"`
	Name     string             `json:"name"`
	NodeType NodeType           `json:"nodeType"`
	Error    string             `json:"error,omitempty"`
	Inputs   map[string]*NodeIO `json:"inputs,omitempty"`
	Outputs  map[string]*NodeIO `json:"outputs,omitempty"`
	Status   NodeStatus         `json:"status,omitempty"`

	Run NodeFunc `json:"-"`
}

type NodeType string

const (
	SourceNode    NodeType = "source"
	ConnectorNode NodeType = "connector"
	EndNode       NodeType = "end"
	DownloadNode  NodeType = "download"
)

type NodeIO struct {
	Id       string   `json:"id"`
	Type     IOType   `json:"type"`
	Label    string   `json:"label"`
	Required bool     `json:"required"`
	ReadOnly bool     `json:"readOnly"`
	Value    any      `json:"value,omitempty"`
	Options  []string `json:"options,omitempty"` // for enums
}

type NodeFunc func(ctx context.Context, node Node, comm chan<- PipelineMessage) error

type IOType string

const (
	IOTypeString    IOType = "string"
	IOTypeNumber    IOType = "number"
	IOTypeBoolean   IOType = "boolean"
	IOTypeDirectory IOType = "directory"
)

type NodeHandle string

const (
	NodeHandleInput  NodeHandle = "input"
	NodeHandleOutput NodeHandle = "output"
)

type NodeStatus string

const (
	StatusPending NodeStatus = "pending"
	StatusRunning NodeStatus = "running"
	StatusSuccess NodeStatus = "success"
	StatusFailed  NodeStatus = "failed"
)
