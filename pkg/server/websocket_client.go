package server

import (
	"github.com/KhoalaS/godel/pkg/pipeline"
	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Id   string
	Conn *websocket.Conn
	Send chan pipeline.PipelineMessage
}
