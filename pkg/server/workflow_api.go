package server

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/KhoalaS/godel/pkg/core"
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/KhoalaS/godel/pkg/registries"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/rs/zerolog/log"
)

type WorkFlowApi struct {
	workflowService core.IWorkflowService
	upgrader        websocket.Upgrader
	clientRegistry  *registries.TypedSyncMap[string, *WebsocketClient]
}

func NewWorkflowApi(workflowService core.IWorkflowService) *WorkFlowApi {
	api := &WorkFlowApi{
		workflowService: workflowService,
		upgrader: websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
		clientRegistry: &registries.TypedSyncMap[string, *WebsocketClient]{},
	}

	api.workflowService.RegisterMessageHandler(api.MessageHandler)
	return api
}

func (api *WorkFlowApi) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /nodes", api.getNodes)
	mux.HandleFunc("POST /pipeline/start", api.startPipeline)
	mux.HandleFunc("/updates/pipeline", api.handlePipelineMessage)

}

func (api *WorkFlowApi) getNodes(w http.ResponseWriter, r *http.Request) {
	nodes := api.workflowService.GetNodes()

	data, err := json.Marshal(nodes)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	w.Write(data)
}

func (api *WorkFlowApi) startPipeline(w http.ResponseWriter, r *http.Request) {
	var gv pipeline.GraphView

	data, err := io.ReadAll(r.Body)
	defer r.Body.Close()

	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	err = json.Unmarshal(data, &gv)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	g := gv.ToPipelineGraph(api.workflowService.GetNodeRegistry())
	p := pipeline.NewPipeline(g, api.workflowService.GetCommChannel())
	api.workflowService.StartPipeline(*p)

	responseData := StartPipelineResponse{
		PipelineId: p.Id,
	}

	responseBytes, err := json.Marshal(responseData)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	w.Write(responseBytes)
}

func (api *WorkFlowApi) handlePipelineMessage(w http.ResponseWriter, r *http.Request) {
	c, err := api.upgrader.Upgrade(w, r, nil)
	if err != nil {
		InternalErrorHandler(w, err)
		return
	}

	clientId := uuid.NewString()
	client := &WebsocketClient{Id: clientId, Conn: c}
	api.clientRegistry.Store(clientId, client)
}

func (api *WorkFlowApi) MessageHandler(message pipeline.PipelineMessage) {
	failedClientIds := []string{}

	for _, client := range api.clientRegistry.All() {
		err := client.Conn.WriteJSON(message)

		if err != nil {
			log.Err(err).Msg("WS write")
			client.Conn.Close()
			failedClientIds = append(failedClientIds, client.Id)
			continue
		}
	}

	for _, id := range failedClientIds {
		api.clientRegistry.Delete(id)
	}
}
