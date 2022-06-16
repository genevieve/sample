package chat

// Ref: https://github.com/gorilla/websocket/blob/master/examples/chat/hub.go

type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Registered rooms
	rooms map[*Room]bool

	// Inbound messages from the clients
	broadcast chan []byte

	// Register requests from the clients
	register chan *Client

	// Unregister requests from the clients
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		rooms:      make(map[*Room]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Register(c *Client) {
	h.register <- c
}

func (h *Hub) Run() {
	for {
		select {

		case client := <-h.register:
			h.clients[client] = true

		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}

		case message := <-h.broadcast:
			for client := range h.clients {
				client.send <- message
			}

		}
	}
}

func (h *Hub) findRoomByID(id string) *Room {
	for room := range h.rooms {
		if room.GetID() == id {
			return room
		}
	}
	return nil
}

func (h *Hub) findRoomByName(name string) *Room {
	for room := range h.rooms {
		if room.Name == name {
			return room
		}
	}
	return nil
}

func (h *Hub) createRoom(name string) *Room {
	room := NewRoom(name)
	go room.Run()
	h.rooms[room] = true
	return room
}

func (h *Hub) findClientByName(name string) *Client {
	for client := range h.clients {
		if client.Name == name {
			return client
		}
	}
	return nil
}
