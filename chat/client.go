package chat

import (
	"ass3/db"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// Initialize client struct with email

// Modify storeMessage method to accept email directly
func (c *client) storeMessage(msg []byte, email string) {
	chatMessage := db.ChatMessage{
		ChatID:    c.room.GetChatID(),
		Sender:    email,
		Message:   string(msg),
		Timestamp: time.Now(),
	}
	err := db.InsertMessage(chatMessage)
	if err != nil {
		log.Println("Error storing message:", err)
	}
}

// Function to get email from cookie
func getEmailFromCookie(r *http.Request) string {
	cookie, err := r.Cookie("email")
	if err != nil {
		log.Println("Error retrieving email from cookie:", err)
		return "" // or handle the error as needed
	}
	return cookie.Value
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{
	ReadBufferSize:  socketBufferSize,
	WriteBufferSize: socketBufferSize,
}

type client struct {
	socket  *websocket.Conn
	receive chan []byte
	room    *room
	email   string // Add a field to store the email of the sender
}

func (c *client) read() {
	defer c.socket.Close()
	for {
		_, msg, err := c.socket.ReadMessage()
		if err != nil {
			return
		}
		// Pass email along with the message to storeMessage
		c.storeMessage(msg, c.email)
		c.room.forward <- msg
	}
}

func (c *client) write() {
	defer c.socket.Close()
	for msg := range c.receive {
		err := c.socket.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			return
		}
	}
}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// Get the email from the cookie
	email := getEmailFromCookie(req)

	socket, err := upgrader.Upgrade(w, req, nil)
	if err != nil {
		log.Println("ServeHTTP:", err)
		http.Error(w, "Failed to upgrade to websocket", http.StatusInternalServerError)
		return
	}
	client := &client{
		socket:  socket,
		receive: make(chan []byte, messageBufferSize),
		room:    r,
		email:   email, // Set the email in the client struct
	}
	r.join <- client
	defer func() { r.leave <- client }()
	go client.write()
	client.read()

}
