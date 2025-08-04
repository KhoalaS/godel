package registries

import "github.com/KhoalaS/godel/pkg/types"

var TransformerRegistry = &TypedSyncMap[string, types.DownloadJobTransformer]{}
