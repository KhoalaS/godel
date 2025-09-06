package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/KhoalaS/godel"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/KhoalaS/godel/pkg/types"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// var upgrader = websocket.Upgrader{}
// for debugging
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var jobs = make(chan *types.DownloadJob, 12)
var comm = make(chan pipeline.PipelineMessage, 96)
var pipelines = make(chan *pipeline.Pipeline, 12)

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
		log.Info().Msg("Starting debug http file server at http://localhost:9999/files/")

		testMux := http.NewServeMux()
		testMux.Handle("/files/", http.StripPrefix("/files/", http.FileServer(http.Dir("./testfiles"))))
		testServer := http.Server{
			Addr:    ":9999",
			Handler: testMux,
		}
		go func() {
			if err := testServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Warn().Err(err).Msg("Listen error on test server")
			}
		}()

		defer testServer.Shutdown(ctx)
	}

	err := godotenv.Load()
	if err != nil {
		log.Warn().Msg("Error loading .env file")
	}

	client := http.Client{}

	pipeline.NodeRegistry["int-input"] = pipeline.CreateIntInputNode()
	pipeline.NodeRegistry["download"] = pipeline.CreateDownloadNode()
	pipeline.NodeRegistry["downloader"] = pipeline.CreateDownloaderNode()
	pipeline.NodeRegistry["basename"] = pipeline.CreateBasenameNode()
	pipeline.NodeRegistry["bytes-input"] = pipeline.CreateBytesInputNode()
	pipeline.NodeRegistry["directory-input"] = pipeline.CreateDirectoryInputNode()

	for i := range *numWorkers {
		wg.Add(1)
		go pipeline.PipelineWorker(ctx, &wg, i, pipelines, &client)
	}

	assetsFS, err := fs.Sub(godel.EmbeddedFiles, "ui/dist/assets")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleRoot)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServerFS(assetsFS)))
	mux.HandleFunc("POST /add", handleAdd)
	mux.HandleFunc("POST /cancel", handleCancel)
	mux.HandleFunc("POST /pause", handlePause)
	mux.HandleFunc("GET /jobs", handleJobs)
	mux.HandleFunc("GET /nodes", handleNodes)
	mux.HandleFunc("POST /pipeline/start", handleStartPipeline)

	mux.HandleFunc("/updates/pipeline", handlePipelineMessage)

	server := &http.Server{
		Addr:    ":9095",
		Handler: corsMiddleWare(mux),
	}

	go func() {
		log.Info().Msg("Starting http server at http://localhost:9095/")
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

func handleRoot(w http.ResponseWriter, r *http.Request) {
	data, err := godel.EmbeddedFiles.ReadFile("ui/dist/index.html")
	if err != nil {
		log.Err(err).Send()
		http.Error(w, "could not read index.html", http.StatusInternalServerError)
		return
	}

	reader := bytes.NewReader(data)
	http.ServeContent(w, r, "index.html", time.Now(), reader)
}

func handleAdd(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)

	var job types.DownloadJob

	err := json.Unmarshal(data, &job)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	job.Id = uuid.New().String()
	job.CancelCh = make(chan struct{}, 1)
	job.PauseCh = make(chan struct{}, 1)
	job.DeleteOnCancel = *deleteOnCancel
	job.Status.Store(types.IDLE)
	if job.Headers == nil {
		job.Headers = map[string]string{}
	}

	if len(job.Transformer) > 0 {
		for _, id := range job.Transformer {
			if tr, ok := registries.TransformerRegistry.Load(id); ok {
				err := tr(&job)
				if err != nil {
					w.WriteHeader(http.StatusBadRequest)
					w.Write(fmt.Appendf(nil, "bad transformer %s", id))
				}
			}
		}
	}

	log.Debug().Any("job", job).Send()

	registries.JobRegistry.Store(job.Id, &job)

	jobs <- &job

	w.Write([]byte(job.Id))
}

func handleCancel(w http.ResponseWriter, r *http.Request) {
	data, _ := io.ReadAll(r.Body)

	jobId := string(data)

	job, ok := registries.JobRegistry.Load(jobId)
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

	job, ok := registries.JobRegistry.Load(jobId)
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

func handleJobs(w http.ResponseWriter, r *http.Request) {
	jobs := registries.JobRegistry.All()

	data, err := json.Marshal(jobs)
	if err != nil {
		responseData, _ := json.Marshal(types.ErrorResponse{
			Error: utils.INTERNAL_ERROR_MESSAGE,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	}

	w.Write(data)
}

func handleNodes(w http.ResponseWriter, r *http.Request) {
	nodes := []pipeline.Node{}

	for _, v := range pipeline.NodeRegistry {
		nodes = append(nodes, v)
	}

	data, err := json.Marshal(nodes)
	if err != nil {
		responseData, _ := json.Marshal(types.ErrorResponse{
			Error: utils.INTERNAL_ERROR_MESSAGE,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	}

	w.Write(data)
}

func handleStartPipeline(w http.ResponseWriter, r *http.Request) {
	var gv pipeline.GraphView

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		responseData, _ := json.Marshal(types.ErrorResponse{
			Error: utils.INTERNAL_ERROR_MESSAGE,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	}

	err = json.Unmarshal(data, &gv)
	if err != nil {
		responseData, _ := json.Marshal(types.ErrorResponse{
			Error: utils.INTERNAL_ERROR_MESSAGE,
		})

		w.WriteHeader(http.StatusInternalServerError)
		w.Write(responseData)
		return
	}

	g := gv.ToPipelineGraph(pipeline.NodeRegistry)
	p := pipeline.NewPipeline(g, comm)
	pipelines <- p
}

func corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func handlePipelineMessage(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Err(err).Msg("WS upgrade")
		return
	}

	client := &pipeline.Client{Conn: c, Send: make(chan pipeline.PipelineMessage, 12)}
	clientId := uuid.NewString()
	pipeline.ClientRegistry.Store(clientId, client)

	defer c.Close()
	defer pipeline.ClientRegistry.Delete(clientId)
	for {
		msg := <-client.Send

		err = c.WriteJSON(msg)
		if err != nil {
			log.Err(err).Str("pipelineId", msg.PipelineId).Str("nodeId", msg.NodeId).Msg("WS write")
			break
		}
	}
}
