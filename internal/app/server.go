package app

import "log"

func (s *ChatServer) Run() {
	for {
		select {
		case client := <-s.Register:
			s.mu.Lock()
			s.Clients[client.UserId] = client
			s.mu.Unlock()
			log.Printf("User %d registered", client.UserId)

		case client := <-s.Unregister:
			s.mu.Lock()
			if _, ok := s.Clients[client.UserId]; ok {
				delete(s.Clients, client.UserId)
				close(client.Send)
			}
			s.mu.Unlock()
			log.Printf("User %d unregistered", client.UserId)
		}
	}
}
