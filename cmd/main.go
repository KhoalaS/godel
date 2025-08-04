package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/KhoalaS/godel/pkg/utils/transformer"
	"github.com/joho/godotenv"
)

var transformerRegistry = map[string]types.DownloadJobTransformer{
	"real-debrid": transformer.RealDebridTransformer,
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jobs := make(chan *types.DownloadJob, 4)

	client := http.Client{}

	numWorkers := 4
	for i := range numWorkers {
		go downloadWorker(ctx, i, jobs, &client)
	}

	job := types.DownloadJob{
		Url:      "http://localhost:8080/files/random.txt",
		Id:       "100",
		Filename: "./testfiles/random_cpy.txt",
		Limit:    1000 * 1024,
	}

	jobs <- &job
	close(jobs)

	<-ctx.Done()
}

func downloadWorker(ctx context.Context, id int, jobs <-chan *types.DownloadJob, client *http.Client) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			for job := range jobs {
				fmt.Printf("Downloading using worker %d\n", id)
				err := utils.Download(client, job, nil)
				if err != nil {
					fmt.Println(err)
				}
			}
		}
	}
}
