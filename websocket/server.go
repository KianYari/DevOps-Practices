package websocket

import "sync"

type Hub struct {
	rooms      map[uint]map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	mu         sync.Mutex
}

func NewHub() *Hub {
	return &Hub{
		rooms:      make(map[uint]map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (hub *Hub) Run() {
	for {
		select {
		case client := <-hub.register:
			hub.RegisterClient(client)
		case client := <-hub.unregister:
			hub.UnregisterClient(client)
		case message := <-hub.broadcast:
			hub.BroadcastMessage(message)
		}
	}
}

func (hub *Hub) RegisterClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	room, ok := hub.rooms[client.roomId]
	if !ok {
		room = make(map[*Client]bool)
		hub.rooms[client.roomId] = room
	}
	room[client] = true
}

func (hub *Hub) UnregisterClient(client *Client) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	room, ok := hub.rooms[client.roomId]
	if ok {
		delete(room, client)
		if len(room) == 0 {
			delete(hub.rooms, client.roomId)
		}
	}
}

func (hub *Hub) BroadcastMessage(message []byte) {
	hub.mu.Lock()
	defer hub.mu.Unlock()

	for _, room := range hub.rooms {
		for client := range room {
			select {
			case client.send <- message:
			default:
				close(client.send)
				delete(room, client)
				if len(room) == 0 {
					delete(hub.rooms, client.roomId)
				}
			}
		}
	}
}
