package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/KhoalaS/godel"
	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/KhoalaS/godel/pkg/utils/transformer"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var jobs = make(chan *types.DownloadJob, 12)
var jobRegistry = &registries.TypedSyncMap[string, *types.DownloadJob]{}
var configs map[string]types.DownloadConfig

func main() {

	var wg sync.WaitGroup

	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	numWorkers := flag.Int("worker", 4, "number of workers")
	flag.Parse()

	log.Info().Msgf("Using %d workers", *numWorkers)

	err := godotenv.Load()
	if err != nil {
		log.Info().Msg("Error loading .env file")
	}

	client := http.Client{}

	configFile, err := os.Open("./configs.json")

	if err == nil {
		configData, err := io.ReadAll(configFile)
		if err != nil {
			log.Fatal().Msg("Could not load configs.json file")
		}

		err = json.Unmarshal(configData, &configs)
		if err != nil {
			log.Warn().Msg("Could not unmarshal configs.json")
			configs = map[string]types.DownloadConfig{}
		}
		configFile.Close()
	} else {
		configs = map[string]types.DownloadConfig{}
	}

	registries.TransformerRegistry.Store("real-debrid", transformer.RealDebridTransformer)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	for i := range *numWorkers {
		wg.Add(1)
		go godel.DownloadWorker(ctx, &wg, i, jobs, &client)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /cancel", handleCancel)

	server := &http.Server{
		Addr:    ":9095",
		Handler: mux,
	}

	log.Info().Msg("Starting http server")

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Listen error")
		}
	}()

	<-ctx.Done()

	// give server 5 seconds to wrap up
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Waiting for workers to exit")
	wg.Wait()

	log.Info().Msg("Server shut down gracefully")
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)

	var job types.DownloadJob

	json.Unmarshal(data, &job)
	job.Id = uuid.New().String()
	job.CancelCh = make(chan struct{})
	job.Status = types.IDLE

	if job.ConfigId != "" {
		if config, exist := configs[job.ConfigId]; exist {
			utils.ApplyConfig(&job, config)
		}
	}

	if len(job.Transformer) > 0 {
		for _, id := range job.Transformer {
			if tr, ok := registries.TransformerRegistry.Load(id); ok {
				err := tr(&job)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(fmt.Appendf(nil, "bad transformer %s", id))
				}
			} else {

			}
		}
	}

	jobRegistry.Store(job.Id, &job)

	jobs <- &job

	w.Write([]byte(job.Id))
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)

	jobId := string(data)

	job, ok := jobRegistry.Load(jobId)
	if !ok {
		log.Info().Msgf("Tried canceling unkown job %s", jobId)
		w.WriteHeader(404)
	} else {
		log.Info().Msgf("Canceled job %s", jobId)
		job.CancelCh <- struct{}{}
	}

}
