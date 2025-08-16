package registries

import "github.com/KhoalaS/godel/pkg/types"

var ConfigReistry = &TypedSyncMap[string, *types.DownloadConfig]{}
