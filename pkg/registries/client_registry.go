package registries

import "github.com/KhoalaS/godel/pkg/types"

var ClientRegistry = &TypedSyncMap[string, *types.Client]{}
