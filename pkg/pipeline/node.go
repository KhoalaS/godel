package pipeline

import (
	"context"
)

type Node struct {
	Id       string             `json:"id,omitempty"`
	Type     string             `json:"type"`
	Name     string             `json:"name"`
	Category NodeCategory       `json:"category,omitempty"`
	Error    string             `json:"error,omitempty"`
	Io       map[string]*NodeIO `json:"io,omitempty"`
	Status   NodeStatus         `json:"status,omitempty"`
	Progress float64            `json:"progress,omitempty"`

	Run NodeFunc `json:"-"`
}

type NodeCategory string

const (
	NodeCategoryInput      NodeCategory = "input"
	NodeCategoryDownloader NodeCategory = "downloader"
	NodeCategoryUtility    NodeCategory = "utility"
)

type NodeIO struct {
	Id          string            `json:"id"`
	ValueType   ValueType         `json:"valueType"`
	Label       string            `json:"label,omitempty"`
	Required    bool              `json:"required"`
	ReadOnly    bool              `json:"readOnly"`
	Value       any               `json:"value,omitempty"`
	Options     []string          `json:"options,omitempty"` // for enums
	Type        IOType            `json:"type"`
	Hooks       map[string]string `json:"hooks,omitempty"`
	HookMapping map[string]string `json:"hookMapping,omitempty"`
}

type NodeFunc func(ctx context.Context, node Node, pipeline IPipeline) error

type ValueType string

const (
	ValueTypeString      ValueType = "string"
	ValueTypeNumber      ValueType = "number"
	ValueTypeBoolean     ValueType = "boolean"
	ValueTypeDirectory   ValueType = "directory"
	ValueTypeDownloadJob ValueType = "downloadjob"
	ValueTypeUnknown     ValueType = "unknown"
)

type IOType string

const (
	IOTypeInput         IOType = "input"
	IOTypeOutput        IOType = "output"
	IOTypePassthrough   IOType = "passthrough"
	IOTypeGenerated     IOType = "generated"
	IOTypeConnectedOnly IOType = "connected_only"
	IOTypeSelection     IOType = "selection"
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
