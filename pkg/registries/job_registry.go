package registries

import "github.com/KhoalaS/godel/pkg/types"

var JobRegistry = &TypedSyncMap[string, *types.DownloadJob]{}
