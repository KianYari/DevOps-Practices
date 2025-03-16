package websocket

import "github.com/gorilla/websocket"

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	roomId uint
	userId uint
}

func NewClient(hub *Hub, conn *websocket.Conn, roomId uint, userId uint) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte),
		roomId: roomId,
		userId: userId,
	}
}

func (client *Client) ReadPump() {
	defer func() {
		client.hub.unregister <- client
		client.conn.Close()
	}()

	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			break
		}
		client.hub.broadcast <- message
	}
}

func (client *Client) WritePump() {
	defer client.conn.Close()
	for message := range client.send {
		err := client.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}
