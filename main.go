package main

import (
	"log"
	"net/http"
	"time"
)

func main() {
	// Create configuration
	config := NewConfig()

	// Initialize stock scraper
	scraper, err := NewStockScraper()
	if err != nil {
		log.Fatal("Failed to initialize scraper:", err)
	}
	defer scraper.Close()

	// Create WebSocket hub
	hub := NewHub()
	go hub.Run()

	// Setup HTTP routes
	http.HandleFunc("/ws", WebSocketHandler(hub))

	// Start WebSocket server
	go func() {
		log.Printf("WebSocket server starting on %s/ws", config.Port)
		if err := http.ListenAndServe(config.Port, nil); err != nil {
			log.Fatal("WebSocket server failed:", err)
		}
	}()

	// Start ticker for regular price updates
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	log.Printf("Starting stock price monitoring for %s every %v", config.StockSymbol, config.Interval)

	for range ticker.C {
		stockPrice, err := scraper.GetStockPrice(config.StockSymbol)
		if err != nil {
			log.Printf("Error getting stock price: %v", err)
			continue
		}
		log.Printf("[%s] %s Stock Price: $%s (Change: %s, %s%%)",
			stockPrice.Timestamp.Format("15:04:05"),
			config.StockSymbol,
			stockPrice.Price,
			stockPrice.ChangeAbsolute,
			stockPrice.ChangePercent)

		// Broadcast to WebSocket clients
		hub.Broadcast(stockPrice)
	}
}
