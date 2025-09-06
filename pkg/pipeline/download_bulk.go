package pipeline

import (
	"context"
	"net/http"
	"net/url"
	"path"
	"time"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

func InferFilename(_url string) (string, error) {
	parsedUrl, err := url.Parse(_url)
	if err != nil {
		return "", err
	}
	name := path.Base(parsedUrl.Path)
	if name == "" {
		name = uuid.NewString()
	}
	return name, nil
}

func DownloadBulk(ctx context.Context, client *http.Client, parent *types.DownloadJob, jobs chan *types.DownloadJob) {

	for _, _url := range parent.Urls {
		childJob := types.NewDownloadJob()
		childJob.DeleteOnCancel = parent.DeleteOnCancel
		childJob.Limit = parent.Limit
		childJob.ParentId = parent.Id
		childJob.Url = _url
		childJob.Headers = parent.Headers
		childJob.IsParent = false
		childJob.DestPath = parent.DestPath
		name, err := InferFilename(_url)
		if err == nil {
			childJob.Filename = name
		}

		registries.JobRegistry.Store(childJob.Id, childJob)
		log.Debug().Str("filename", childJob.Filename).Str("url", childJob.Url).Msg("Add child job to jobs channel")

	SelectLoop:
		for {
			select {
			case <-ctx.Done():
				return
			case jobs <- childJob:
				break SelectLoop
			default:
				log.Warn().Msg("jobs channel full waiting...")
				time.Sleep(time.Second * 1)
			}
		}
	}
}
