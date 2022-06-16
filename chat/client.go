package chat

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// Ref: https://github.com/gorilla/websocket/blob/master/examples/chat/client.go

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'\n'}
)

type Client struct {
	Name string    `json:"name"`
	ID   uuid.UUID `json:"id"`

	logger *zap.Logger
	conn   *websocket.Conn

	hub   *Hub
	rooms map[*Room]bool

	send chan []byte
}

func NewClient(name string, conn *websocket.Conn, hub *Hub, logger *zap.Logger) *Client {
	id := uuid.New()
	return &Client{
		ID:     id,
		Name:   name,
		logger: logger.With(zap.String("client_id", id.String())),
		conn:   conn,
		hub:    hub,
		send:   make(chan []byte, 256),
		rooms:  make(map[*Room]bool),
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) ReadPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)

	_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))

	c.conn.SetPongHandler(func(string) error {
		_ = c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, jsonMessage, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.logger.Error(err.Error())
			}
			break
		}

		c.handleNewMessage(jsonMessage)
	}
}

const SendMessageAction = "send-message"
const JoinRoomAction = "join-room"
const LeaveRoomAction = "leave-room"
const RoomJoinedAction = "room-joined"

func (c *Client) handleNewMessage(jsonMessage []byte) {
	var message Message
	if err := json.Unmarshal(jsonMessage, &message); err != nil {
		c.logger.Error(err.Error())
		return
	}

	// Attach the client object as the sender of the messsage.
	message.Sender = c

	switch message.Action {
	case SendMessageAction:
		if room := c.hub.findRoomByID(message.Target.GetID()); room != nil {
			room.broadcast <- &message
		}

	case LeaveRoomAction:
		room := c.hub.findRoomByID(message.Message)
		if room == nil {
			c.logger.Error("attempted to leave room that does not exist")
			return
		}

		// TODO: Since a room is 1:1, if one person leaves,
		// the other person should be notified
		delete(c.rooms, room)

		room.unregister <- c

	case JoinRoomAction:
		// TODO: Find client by ID
		target := c.hub.findClientByName(message.Message)
		if target == nil {
			c.logger.Error("attempted to join room with target that does not exist")
			return
		}

		roomName := message.Message
		// TODO: Create room name from the initiator ID + target ID
		// roomName := message.Message + client.ID.String()

		c.joinRoom(roomName, target)
		target.joinRoom(roomName, c)
	}
}

func (c *Client) joinRoom(roomName string, sender *Client) {
	room := c.hub.findRoomByName(roomName)
	if room == nil {
		room = c.hub.createRoom(roomName)
	}

	if sender == nil {
		return
	}

	if !c.isInRoom(room) {
		c.rooms[room] = true
		room.register <- c

		message := Message{
			Action: RoomJoinedAction,
			Target: room,
			Sender: sender,
		}

		c.send <- message.encode()
	}
}

func (c *Client) isInRoom(room *Room) bool {
	if _, ok := c.rooms[room]; ok {
		return true
	}

	return false
}

// writePump pumps messages from the hub to the websocket connection.
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error(err.Error())
				return
			}
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message.
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				c.logger.Error(err.Error())
				return
			}
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
