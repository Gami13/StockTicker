package main

import (
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

		hub.register <- conn

		// Keep the connection alive and handle disconnection
		go func() {
			defer func() {
				hub.unregister <- conn
			}()

			for {
				_, _, err := conn.ReadMessage()
				if err != nil {
					break
				}
			}
		}()
	}
}
