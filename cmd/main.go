package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/KhoalaS/godel/pkg/utils/transformer"
	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)

	defer stop()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var configs map[string]types.DownloadConfig
	configFile, err := os.Open("./configs.json")

	if err == nil {
		configData, err := io.ReadAll(configFile)
		if err != nil {
			log.Fatal("Error loading configs.json file")
		}

		json.Unmarshal(configData, &configs)
		configFile.Close()
	} else {
		configs = map[string]types.DownloadConfig{}
	}

	registries.TransformerRegistry.Store("real-debrid", transformer.RealDebridTransformer)

	jobs := make(chan *types.DownloadJob, 4)

	client := http.Client{}

	numWorkers := 4
	for i := range numWorkers {
		go downloadWorker(ctx, i, jobs, &client)
	}

	job := types.DownloadJob{
		Url:      "http://localhost:8080/files/stuff.zip",
		Id:       "100",
		Filename: "./testfiles/stuff_cpy.zip",
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
