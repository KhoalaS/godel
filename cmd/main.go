package main

import (
	"fmt"
	"log"
	"net/http"

	//"strings"
	"sync"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/KhoalaS/godel/pkg/utils/transformer"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jobs := make(chan *types.DownloadJob, 4)

	client := http.Client{}

	var wg sync.WaitGroup

	numWorkers := 4
	for i := range numWorkers {
		wg.Add(1)
		go downloadWorker(i, jobs, &wg, &client)
	}

	rdJob, err := transformer.RealDebridTransformer(types.DownloadJob{
		Url:      "https://rapidgator.net/file/0347f6bddc27aa6dcf7c64df751ae783",
		Id:       100,
		Filename: "",
		Password: "",
	})

	if err != nil {
		log.Fatal(err)
	}

	jobs <- &rdJob

	close(jobs)

	wg.Wait()
}

func downloadWorker(id int, jobs <-chan *types.DownloadJob, wg *sync.WaitGroup, client *http.Client) {
	for job := range jobs {
		fmt.Printf("Downloading using worker %d\n", id)
		err := utils.Download(client, job.Url, job.Filename, nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	defer wg.Done()
}
