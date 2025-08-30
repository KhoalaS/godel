package pipeline

import (
	"context"
)

type Node struct {
	Id       string             `json:"id,omitempty"`
	Type     string             `json:"type"`
	Name     string             `json:"name"`
	NodeType NodeType           `json:"nodeType,omitempty"`
	Error    string             `json:"error,omitempty"`
	Io       map[string]*NodeIO `json:"io,omitempty"`
	Status   NodeStatus         `json:"status,omitempty"`

	Run NodeFunc `json:"-"`
}

type NodeType string

const (
	InputNode      NodeType = "input"
	DownloaderNode NodeType = "downloader"
)

type NodeIO struct {
	Id        string            `json:"id"`
	ValueType ValueType         `json:"valueType"`
	Label     string            `json:"label"`
	Required  bool              `json:"required"`
	ReadOnly  bool              `json:"readOnly"`
	Value     any               `json:"value,omitempty"`
	Options   []string          `json:"options,omitempty"` // for enums
	Type      IOType            `json:"type"`
	Hooks     map[string]string `json:"hooks,omitempty"`
}

type NodeFunc func(ctx context.Context, node Node, comm chan<- PipelineMessage) error

type ValueType string

const (
	ValueTypeString    ValueType = "string"
	ValueTypeNumber    ValueType = "number"
	ValueTypeBoolean   ValueType = "boolean"
	ValueTypeDirectory ValueType = "directory"
)

type IOType string

const (
	IOTypeInput       IOType = "input"
	IOTypeOutput      IOType = "output"
	IOTypePassthrough IOType = "passthrough"
	IOTypeGenerated   IOType = "generated"
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
