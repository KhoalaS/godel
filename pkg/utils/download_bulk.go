package utils

import (
	"context"
	"net/http"
	"path/filepath"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/rs/zerolog/log"
)

func DownloadBulk(ctx context.Context, client *http.Client, parent *types.DownloadJob, jobs chan *types.DownloadJob) {

	for _, _url := range parent.Urls {
		childJob := types.NewDownloadJob()
		childJob.DeleteOnCancel = parent.DeleteOnCancel
		childJob.Limit = parent.Limit
		childJob.ParentId = parent.Id
		childJob.Transformer = parent.Transformer
		childJob.Url = _url
		childJob.Headers = parent.Headers
		childJob.IsParent = false
		name, err := InferFilename(_url)
		if err == nil {
			childJob.Filename = filepath.Join(parent.Filename, name)
		}

		trError := false
		if len(parent.Transformer) > 0 {
			for _, id := range parent.Transformer {
				if tr, ok := registries.TransformerRegistry.Load(id); ok {
					err := tr(childJob)
					if err != nil {
						log.Warn().Str("url", childJob.Url).Msg("Error applying transformer in child of bulk download")
						trError = true
						break
					}
				}
			}
		}

		if trError {
			continue
		}

		registries.JobRegistry.Store(childJob.Id, childJob)
		log.Debug().Str("filename", childJob.Filename).Str("url", childJob.Url).Msg("Add child job to jobs channel")
		jobs <- childJob
	}

}
