package main

import (
	"net/http"

	"github.com/ayananerv/do-ur-chat/internal/app"
)

func main() {
	server := app.NewChatServer()
	go server.Run()

	http.HandleFunc("/ws", server.HandleWebSocket)

	println("Server started on port 8080")
	http.ListenAndServe(":8080", nil)
}
