package server

import (
	"github.com/gorilla/websocket"
)

type WebsocketClient struct {
	Id   string
	Conn *websocket.Conn
}
