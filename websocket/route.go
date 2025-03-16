package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func SetupWSRoutes(router *gin.Engine) {
	hub := NewHub()
	go hub.Run()

	wsMiddleware := NewWebSocketMiddleware()
	router.GET("/ws", wsMiddleware.Upgrade, func(c *gin.Context) {
		conn, exists := c.Get("wsConn")
		if !exists {
			c.JSON(500, gin.H{"error": "Failed to get WebSocket connection"})
			return
		}
		client := NewClient(hub, conn.(*websocket.Conn), 0, 0)
		hub.register <- client
		go client.ReadPump()
		go client.WritePump()
	})
}
