package auth

import "github.com/KhoalaS/godel/pkg/types"

type AuthProvider func(job *types.DownloadJob) (types.Credentials, error)
