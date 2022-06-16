package chat

import (
	"github.com/google/uuid"
)

type Room struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan *Message
}

func NewRoom(name string) *Room {
	return &Room{
		ID:         uuid.New(),
		Name:       name,
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (r *Room) Run() {
	for {
		select {
		case client := <-r.register:
			r.registerClient(client)
		case client := <-r.unregister:
			r.unregisterClient(client)
		case message := <-r.broadcast:
			r.broadcastToClients(message)
		}
	}
}

// broadcastToClients will broadcast the message to clients
// in the room with the exception of the original sender
func (r *Room) broadcastToClients(message *Message) {
	msg := message.encode()
	for client := range r.clients {
		if client.ID != message.Sender.ID {
			client.send <- msg
		}
	}
}

// registerClient will add the client to the list of clients in the room.
func (r *Room) registerClient(client *Client) {
	r.clients[client] = true
}

// unregisterClient will remove the client to the list of clients in the room.
func (r *Room) unregisterClient(client *Client) {
	delete(r.clients, client)
}

func (r *Room) GetID() string {
	return r.ID.String()
}
