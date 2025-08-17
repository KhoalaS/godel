package registries

import "github.com/KhoalaS/godel/pkg/types"

var ConfigRegistry = &TypedSyncMap[string, *types.DownloadConfig]{}
