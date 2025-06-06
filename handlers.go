package main

import (
	_ "embed"
	"log"
	"net/http"
)

// WebSocketHandler handles WebSocket connections
func WebSocketHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Printf("WebSocket upgrade error: %v", err)
			return
		}

		client := &Client{
			Conn:   conn,
			Symbol: "", // Will be set when client sends a request
		}

		// Keep the connection alive and handle messages
		go func() {
			defer func() {
				hub.unregister <- client
			}()

			for {
				var message ClientMessage
				err := conn.ReadJSON(&message)
				if err != nil {
					log.Printf("WebSocket read error: %v", err)
					break
				}

				// Handle stock symbol request
				if message.Type == "subscribe" && message.Symbol != "" {
					client.Symbol = message.Symbol
					hub.register <- client
					log.Printf("Client subscribed to %s", message.Symbol)
				}
			}
		}()
	}
}

//go:embed index.html
var html string

// StaticHandler serves the HTML page
func StaticHandler() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, html)
	}
}
