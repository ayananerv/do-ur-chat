package app

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn   *websocket.Conn
	Send   chan []byte // buffer queue
	UserId int64
}

// ChatServer (Hub) to manage all connections
type ChatServer struct {
	// Read Write Lock: to protect `Clients`
	mu sync.RWMutex
	// Core dict: map UserId -> Client
	Clients map[int64]*Client
	// channels handling register / deactivate
	Register   chan *Client
	Unregister chan *Client
}

// initialization
func NewChatServer() *ChatServer {
	return &ChatServer{
		Clients:    make(map[int64]*Client),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}
