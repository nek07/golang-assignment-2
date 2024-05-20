// room.go
package chat

import (
	"net/http"

	"github.com/google/uuid"
)

type room struct {
	chatID  string
	forward chan []byte
	join    chan *client
	leave   chan *client
	clients map[*client]bool
}

func NewRoom(chatID string) *room {
	return &room{

		chatID:  uuid.New().String(),
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
	}
}

func (r *room) Run() {
	for {
		select {
		case client := <-r.join:
			r.clients[client] = true
		case client := <-r.leave:
			delete(r.clients, client)
			close(client.receive)
		case msg := <-r.forward:
			for client := range r.clients {
				select {
				case client.receive <- msg:
				default:
					delete(r.clients, client)
					close(client.receive)
				}
			}
		}
	}
}
func (r *room) HandleRoom(w http.ResponseWriter, req *http.Request) {

	r.ServeHTTP(w, req)
}

func (r *room) GetChatID() string {
	return r.chatID
}
