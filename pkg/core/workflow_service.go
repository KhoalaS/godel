package core

import (
	"context"
	"net/http"
	"os"
	"sync"

	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type IWorkflowService interface {
	GetNodes() []*pipeline.Node
	GetNodeRegistry() map[string]*pipeline.Node
	GetCommChannel() chan pipeline.PipelineMessage
	StartPipeline(pipeline pipeline.Pipeline)
	HandlePipelineMessage(message pipeline.PipelineMessage)
	RegisterNode(id string, nodeConstructor func() pipeline.Node)
	RegisterMessageHandler(handler func(message pipeline.PipelineMessage))
}

type WorkflowService struct {
	comm            chan pipeline.PipelineMessage
	pipelines       chan *pipeline.Pipeline
	deleteOnCancel  bool
	debugMode       bool
	numWorkers      int
	nodeRegistry    map[string]*pipeline.Node
	logger          zerolog.Logger
	wg              *sync.WaitGroup
	client          *http.Client
	messageHandlers map[string]func(messsage pipeline.PipelineMessage)
}

type WorkflowServiceConfig struct {
	CommChannelSize     int
	PipelineChannelSize int
	DeleteOnCancel      bool
	DebugMode           bool
	NumWorkers          int
	LogLevel            zerolog.Level
}

func NewWorkflowService(ctx context.Context, config WorkflowServiceConfig) *WorkflowService {
	commChannelSize := config.CommChannelSize
	if commChannelSize == 0 {
		commChannelSize = 96
	}

	pipelineChannelSize := config.PipelineChannelSize
	if pipelineChannelSize == 0 {
		pipelineChannelSize = 12
	}

	numWorkers := config.NumWorkers
	if numWorkers == 0 {
		numWorkers = 4
	}

	service := &WorkflowService{
		comm:            make(chan pipeline.PipelineMessage, commChannelSize),
		pipelines:       make(chan *pipeline.Pipeline, pipelineChannelSize),
		deleteOnCancel:  config.DeleteOnCancel,
		debugMode:       config.DebugMode,
		numWorkers:      numWorkers,
		nodeRegistry:    make(map[string]*pipeline.Node),
		logger:          log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(config.LogLevel),
		messageHandlers: make(map[string]func(messsage pipeline.PipelineMessage)),
		wg:              &sync.WaitGroup{},
		client:          &http.Client{},
	}

	for i := range numWorkers {
		service.wg.Add(1)
		go pipeline.PipelineWorker(ctx, service.wg, i, service.pipelines, service.client)
	}

	go func(service *WorkflowService) {
		for {
			message := <-service.comm
			service.HandlePipelineMessage(message)
		}
	}(service)

	return service
}

func (s *WorkflowService) GetNodes() []*pipeline.Node {
	nodes := []*pipeline.Node{}
	for _, node := range s.nodeRegistry {
		nodes = append(nodes, node)
	}

	return nodes
}

func (s *WorkflowService) RegisterNode(id string, nodeConstructor func() pipeline.Node) {
	node := nodeConstructor()
	if _, ex := s.nodeRegistry[id]; ex {
		s.logger.Warn().Str("id", id).Msg("tried registering a node with existing id")
		return
	}
	s.nodeRegistry[id] = &node
}

func (s *WorkflowService) StartPipeline(pipeline pipeline.Pipeline) {
	pipeline.Comm = s.comm
	s.pipelines <- &pipeline
}

func (s *WorkflowService) HandlePipelineMessage(message pipeline.PipelineMessage) {
	s.logger.Info().Str("level", string(message.Level)).Str("type", string(message.Type)).Any("data", message.Data).Msg("incoming message")
	for _, handler := range s.messageHandlers {
		handler(message)
	}
}

func (s *WorkflowService) RegisterMessageHandler(handler func(message pipeline.PipelineMessage)) {
	handlerId := uuid.NewString()
	s.messageHandlers[handlerId] = handler
}

func (s *WorkflowService) GetCommChannel() chan pipeline.PipelineMessage {
	return s.comm
}

func (s *WorkflowService) GetNodeRegistry() map[string]*pipeline.Node {
	return s.nodeRegistry
}
