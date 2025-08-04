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

var transformerRegistry = map[string]types.DownloadJobTransformer{
	"real-debrid": transformer.RealDebridTransformer,
}

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

	job := types.DownloadJob{
		Url:      "http://localhost:8080/files/random.txt",
		Id:       "100",
		Filename: "./testfiles/random_cpy.txt",
		Limit:    1000 * 1024,
	}

	jobs <- &job

	close(jobs)

	wg.Wait()
}

func downloadWorker(id int, jobs <-chan *types.DownloadJob, wg *sync.WaitGroup, client *http.Client) {
	for job := range jobs {
		fmt.Printf("Downloading using worker %d\n", id)
		err := utils.Download(client, job, nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	defer wg.Done()
}
