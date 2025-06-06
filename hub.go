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
	clients    map[*Client]bool
	broadcast  chan StockPrice
	register   chan *Client
	unregister chan *Client
	mutex      sync.RWMutex
}

// NewHub creates a new WebSocket hub
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan StockPrice),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub and handles client connections and broadcasts
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.registerClient(client)

		case client := <-h.unregister:
			h.unregisterClient(client)

		case stockPrice := <-h.broadcast:
			h.broadcastToClients(stockPrice)
		}
	}
}

func (h *Hub) registerClient(client *Client) {
	h.mutex.Lock()
	h.clients[client] = true
	h.mutex.Unlock()
	log.Printf("Client connected requesting %s. Total clients: %d", client.Symbol, len(h.clients))
}

func (h *Hub) unregisterClient(client *Client) {
	h.mutex.Lock()
	if _, ok := h.clients[client]; ok {
		delete(h.clients, client)
		client.Conn.Close()
	}
	h.mutex.Unlock()
	log.Printf("Client disconnected. Total clients: %d", len(h.clients))
}

func (h *Hub) broadcastToClients(stockPrice StockPrice) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.clients {
		// Only send to clients requesting this specific stock symbol
		if client.Symbol == stockPrice.Symbol {
			if err := client.Conn.WriteJSON(stockPrice); err != nil {
				log.Printf("Error writing to client: %v", err)
				delete(h.clients, client)
				client.Conn.Close()
			}
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

// GetRequestedSymbols returns a slice of unique stock symbols requested by clients
func (h *Hub) GetRequestedSymbols() []string {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	symbolMap := make(map[string]bool)
	for client := range h.clients {
		if client.Symbol != "" {
			symbolMap[client.Symbol] = true
		}
	}

	symbols := make([]string, 0, len(symbolMap))
	for symbol := range symbolMap {
		symbols = append(symbols, symbol)
	}
	return symbols
}
