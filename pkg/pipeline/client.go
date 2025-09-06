package pipeline

import (
	"github.com/gorilla/websocket"
)

type Client struct {
	Conn *websocket.Conn
	Send chan PipelineMessage
}
