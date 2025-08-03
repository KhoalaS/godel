package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
)

func main() {
	jobs := make(chan *types.DownloadJob, 4)

	client := http.Client{
		Timeout: 15 * time.Second,
	}

	var wg sync.WaitGroup

	numWorkers := 4
	for i := range numWorkers {
		wg.Add(1)
		go downloadWorker(i, jobs, &wg, &client)
	}

	urls := []string{
		"http://localhost:8080/files/test.txt",
		"http://localhost:8080/files/test.txt",
		"http://localhost:8080/files/test.txt",
		"http://localhost:8080/files/test.txt",
		"http://localhost:8080/files/video.mp4",
		"http://localhost:8080/files/random.txt",
	}

	for idx, url := range urls {
		spl := strings.Split(url, "/")
		filename := spl[len(spl)-1]
		jobs <- &types.DownloadJob{
			Url:      url,
			Id:       idx,
			Filename: fmt.Sprintf("./testfiles/%d_%s", idx, filename),
		}
	}

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
