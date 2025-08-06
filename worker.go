package godel

import (
	"context"
	"fmt"
	"net/http"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func DownloadWorker(ctx context.Context, id int, jobs <-chan *types.DownloadJob, client *http.Client) {
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("Done signal send for worker %d\n", id)
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			fmt.Printf("Downloading using worker %d\n", id)
			err := utils.Download(ctx, client, job, nil)
			if err != nil {
				fmt.Println("error during download", err)
			}

		}
	}
}
