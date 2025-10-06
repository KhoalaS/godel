package server

import (
	"bytes"
	"context"
	"flag"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/KhoalaS/godel"
	"github.com/KhoalaS/godel/pkg/core"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var deleteOnCancel *bool
var debugMode *bool

func RunServer() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

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

	workflowService := core.NewWorkflowService(ctx, core.WorkflowServiceConfig{
		NumWorkers:     *numWorkers,
		DeleteOnCancel: *deleteOnCancel,
		DebugMode:      *debugMode,
		LogLevel:       logLevel,
	})

	workflowService.RegisterNode("int-input", pipeline.CreateIntInputNode)
	workflowService.RegisterNode("download", pipeline.CreateDownloadNode)
	workflowService.RegisterNode("downloader", pipeline.CreateDownloaderNode)
	workflowService.RegisterNode("basename", pipeline.CreateBasenameNode)
	workflowService.RegisterNode("bytes-input", pipeline.CreateBytesInputNode)
	workflowService.RegisterNode("directory-input", pipeline.CreateDirectoryInputNode)
	workflowService.RegisterNode("rd-downloader", pipeline.CreateRealdebridNode)
	workflowService.RegisterNode("suffix", pipeline.CreateSuffixNode)
	workflowService.RegisterNode("http-request", pipeline.CreateHttpRequestNode)
	workflowService.RegisterNode("display", pipeline.CreateDisplayNode)
	workflowService.RegisterNode("transmission-service", pipeline.CreateTransmissionNode)
	workflowService.RegisterNode("add-torrent", pipeline.CreateAddTorrentNode)
	workflowService.RegisterNode("pixeldrain", pipeline.CreatePixeldrainNode)
	workflowService.RegisterNode("unrar", pipeline.CreateUnrarNode)

	assetsFS, err := fs.Sub(godel.EmbeddedFiles, "ui/dist/assets")
	if err != nil {
		log.Fatal().Err(err).Send()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", handleRoot)
	mux.Handle("GET /assets/", http.StripPrefix("/assets/", http.FileServerFS(assetsFS)))

	api := NewWorkflowApi(workflowService)
	api.RegisterRoutes(mux)

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

	log.Info().Msg("Server shut down gracefully")
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	data, err := godel.EmbeddedFiles.ReadFile("ui/dist/index.html")
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	reader := bytes.NewReader(data)
	http.ServeContent(w, r, "index.html", time.Now(), reader)
}

func corsMiddleWare(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func InternalErrorHandler(w http.ResponseWriter, err error) {
	log.Error().Err(err).Msg("Internal server error")
	http.Error(w, utils.INTERNAL_ERROR_MESSAGE, http.StatusInternalServerError)
}
