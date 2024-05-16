// room.go
package chat

import "net/http"

type room struct {
	forward chan []byte
	join    chan *client1
	leave   chan *client1
	clients map[*client1]bool
}

func NewRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client1),
		leave:   make(chan *client1),
		clients: make(map[*client1]bool),
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
