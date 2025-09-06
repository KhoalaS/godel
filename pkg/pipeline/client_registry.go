package pipeline

import (
	"github.com/KhoalaS/godel/pkg/registries"
)

var ClientRegistry = &registries.TypedSyncMap[string, *Client]{}
