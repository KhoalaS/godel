package registries

import "github.com/KhoalaS/godel/pkg/types"

var AuthRegistry = &TypedSyncMap[string, types.Credentials]{}
