package utils

import (
	"context"
	"net/http"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
)

func DownloadBulk(ctx context.Context, client *http.Client, urls []string, job *types.DownloadJob, headers map[string]string, jobs chan<- *types.DownloadJob) {

	parentJob := types.NewDownloadJob()
	parentJob.BytesDownloaded = 0
	parentJob.Size = len(urls)
	parentJob.IsParent = true

	registries.JobRegistry.Store(parentJob.Id, parentJob)

	jobs <- parentJob

	for _, _url := range urls {
		childJob := types.NewDownloadJob()
		childJob.DeleteOnCancel = job.DeleteOnCancel
		childJob.Limit = job.Limit
		childJob.ParentId = parentJob.Id
		childJob.Transformer = job.Transformer
		childJob.Url = _url

		registries.JobRegistry.Store(childJob.Id, childJob)
		jobs <- childJob
	}

}
