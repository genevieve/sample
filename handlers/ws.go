package handlers

import (
	"net/http"

	"github.com/genevieve/sample/chat"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

type Chat struct {
	hub    *chat.Hub
	logger *zap.Logger
}

func NewChat(hub *chat.Hub, logger *zap.Logger) *Chat {
	return &Chat{
		hub:    hub,
		logger: logger,
	}
}

func (h Chat) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name, ok := r.URL.Query()["name"]
	if !ok {
		h.logger.Error("invalid query `name`")
		http.Error(w, "invalid query `name`", http.StatusBadRequest)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.logger.Error(err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	client := chat.NewClient(name[0], conn, h.hub, h.logger)
	go client.WritePump()
	go client.ReadPump()

	h.hub.Register(client)
}
