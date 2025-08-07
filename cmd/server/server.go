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

var deleteOnCancel *bool
var debugMode *bool

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	var wg sync.WaitGroup

	numWorkers := flag.Int("worker", 4, "number of workers")
	deleteOnCancel = flag.Bool("del-cancel", true, "wether to delete files of canceled downloads")
	debugMode = flag.Bool("debug", false, "runs the server in debug mode, starts a fileserver on port 8080 and serves the ./testfiles directory")

	flag.Parse()

	logLevel := zerolog.InfoLevel
	if *debugMode {
		logLevel = zerolog.DebugLevel
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(logLevel)

	log.Info().Msgf("Using %d workers", *numWorkers)

	if *deleteOnCancel {
		log.Info().Msg("Deleting files on cancel")
	}

	if *debugMode {
		log.Info().Msg("Starting debug http file server")

		testMux := http.NewServeMux()
		testMux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./testfiles"))))
		testServer := http.Server{
			Addr:    ":8080",
			Handler: testMux,
		}
		go func() {
			if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Err(err).Msg("Listen error on test server")
			}
		}()

		defer testServer.Shutdown(ctx)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("Error loading .env file")
	}

	client := http.Client{}

	loadConfig()

	registries.TransformerRegistry.Store("real-debrid", transformer.RealDebridTransformer)

	for i := range *numWorkers {
		wg.Add(1)
		go godel.DownloadWorker(ctx, &wg, i, jobs, &client)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /cancel", handleCancel)
	mux.HandleFunc("POST /pause", handlePause)

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
	job.CancelCh = make(chan struct{}, 1)
	job.PauseCh = make(chan struct{}, 1)
	job.DeleteOnCancel = *deleteOnCancel
	job.Status.Store(types.IDLE)

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
		log.Info().Str("id", jobId).Msg("Tried canceling unknown job")
		w.WriteHeader(404)
	} else {
		currentState := job.Status.Load()
		switch currentState {
		case types.PAUSED:
			log.Info().Str("filename", job.Filename).Str("id", jobId).Msg("Canceling paused job")
			job.Status.Store(types.CANCELED)
			if job.DeleteOnCancel {
				go func() {
					os.Remove(job.Filename)
				}()
			}
		case types.DOWNLOADING:
			select {
			case job.CancelCh <- struct{}{}:
			default:
				log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("Cancel signal already sent")
			}
			log.Info().Str("filename", job.Filename).Str("id", jobId).Msg("Canceled job")
		default:
			log.Warn().Str("id", jobId).Str("status", string(currentState.(types.DownloadState))).Msg("Job could not be canceled in current state")
		}

	}
}

func handlePause(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)

	jobId := string(data)

	job, ok := jobRegistry.Load(jobId)
	if !ok {
		log.Info().Str("id", jobId).Msg("Tried pausing unkown job")
		w.WriteHeader(404)
	} else {
		currentState := job.Status.Load()
		switch currentState {
		case types.PAUSED:
			// discard old pause signal
			select {
			case <-job.PauseCh:
			default:
			}
			jobs <- job
			log.Info().Str("id", jobId).Msg("Resumed job")
		case types.DOWNLOADING:
			select {
			case job.PauseCh <- struct{}{}:
			default:
				log.Warn().Str("filename", job.Filename).Str("id", job.Id).Msg("Pause signal already sent")
			}
			log.Info().Str("id", jobId).Msg("Paused job")
		default:
			log.Warn().Str("id", jobId).Str("status", string(currentState.(types.DownloadState))).Msg("Job could not be paused/resumed in current state")
		}
	}
}

func loadConfig() {
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
}
