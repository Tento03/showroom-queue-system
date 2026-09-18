package controllers

import (
	"backend-queue/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow semua origin untuk dev
	},
}

type WebSocketController struct{}

func NewWebSocketController() *WebSocketController {
	return &WebSocketController{}
}

func (ctrl *WebSocketController) Handle(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upgrade connection"})
		return
	}

	services.Hub.Register(conn)

	// Jaga koneksi tetap hidup
	go func() {
		defer func() {
			services.Hub.Unregister(conn)
		}()
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				break
			}
		}
	}()
}
