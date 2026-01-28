package app

import (
	"log"
	"net/http"
	"strconv"

	"github.com/ayananerv/do-ur-chat/api/gen/message"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

var upgrader = websocket.Upgrader{
	// dev
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *ChatServer) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// ?uid=xxx
	uidStr := r.URL.Query().Get("uid")
	uid, _ := strconv.ParseInt(uidStr, 10, 64)

	// upgrade protocol
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Upgrade error: ", err)
	}

	// create Client object
	client := &Client{Conn: conn, Send: make(chan []byte, 256), UserId: uid}

	// register to HUB
	s.Register <- client

	go func() {
		defer func() {
			conn.Close()
		}()
		for {
			select {
			case msg, ok := <-client.Send:
				if !ok {
					return
				}
				conn.WriteMessage(websocket.BinaryMessage, msg)
			}
		}
	}()

	// ReadPump
	for {
		_, data, err := conn.ReadMessage()
		if err != nil {
			s.Unregister <- client
			break
		}

		msg := &message.ChatMessage{}
		if err := proto.Unmarshal(data, msg); err != nil {
			log.Println("Proto unmarshal error: ", err)
			continue
		}

		log.Println("receive message: %s from %d to %d", msg.Content, msg.SenderId, msg.ReceiverId)

		s.mu.RLock()
		targetClient, ok := s.Clients[msg.ReceiverId]
		s.mu.RUnlock()

		if ok {
			// receiver is online, we can send it directly
			targetClient.Send <- data
		} else {
			log.Printf("User %d is offline, the message is stored in DB", msg.ReceiverId)
		}
	}
}
