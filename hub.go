package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow connections from any origin
	},
}

// Hub maintains the set of active clients and broadcasts messages to the clients
type Hub struct {
	clients    map[*websocket.Conn]bool
	broadcast  chan StockPrice
	register   chan *websocket.Conn
	unregister chan *websocket.Conn
	mutex      sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*websocket.Conn]bool),
		broadcast:  make(chan StockPrice),
		register:   make(chan *websocket.Conn),
		unregister: make(chan *websocket.Conn),
	}
}

// Run starts the hub and handles client connections and broadcasts
func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.registerClient(conn)

		case conn := <-h.unregister:
			h.unregisterClient(conn)

		case stockPrice := <-h.broadcast:
			h.broadcastToClients(stockPrice)
		}
	}
}

func (h *Hub) registerClient(conn *websocket.Conn) {
	h.mutex.Lock()
	h.clients[conn] = true
	h.mutex.Unlock()
	log.Printf("Client connected. Total clients: %d", len(h.clients))
}

func (h *Hub) unregisterClient(conn *websocket.Conn) {
	h.mutex.Lock()
	if _, ok := h.clients[conn]; ok {
		delete(h.clients, conn)
		conn.Close()
	}
	h.mutex.Unlock()
	log.Printf("Client disconnected. Total clients: %d", len(h.clients))
}

func (h *Hub) broadcastToClients(stockPrice StockPrice) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for conn := range h.clients {
		if err := conn.WriteJSON(stockPrice); err != nil {
			log.Printf("Error writing to client: %v", err)
			delete(h.clients, conn)
			conn.Close()
		}
	}
}

// Broadcast sends a stock price update to all connected clients
func (h *Hub) Broadcast(stockPrice StockPrice) {
	select {
	case h.broadcast <- stockPrice:
	default:
		log.Println("Broadcast channel full, dropping message")
	}
}
