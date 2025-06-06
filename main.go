package main

import (
	"log"
	"net/http"
	"sync"
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
	http.HandleFunc("/", StaticHandler())

	// Start HTTP server
	go func() {
		log.Printf("HTTP server starting on %s", config.Port)
		if err := http.ListenAndServe(config.Port, nil); err != nil {
			log.Fatal("HTTP server failed:", err)
		}
	}()

	// Track active stock monitoring goroutines
	activeStocks := make(map[string]bool)
	var stockMutex sync.RWMutex

	// Ticker for checking which stocks to monitor
	ticker := time.NewTicker(config.Interval)
	defer ticker.Stop()

	log.Printf("Stock ticker service started. Monitoring interval: %v", config.Interval)

	for range ticker.C {
		requestedSymbols := hub.GetRequestedSymbols()

		stockMutex.Lock()
		// Start monitoring new symbols
		for _, symbol := range requestedSymbols {
			if !activeStocks[symbol] {
				activeStocks[symbol] = true
				go monitorStock(symbol, scraper, hub, &activeStocks, &stockMutex)
			}
		}

		// Clean up symbols that are no longer requested
		currentSymbols := make(map[string]bool)
		for _, symbol := range requestedSymbols {
			currentSymbols[symbol] = true
		}

		for symbol := range activeStocks {
			if !currentSymbols[symbol] {
				activeStocks[symbol] = false // Mark for cleanup
			}
		}
		stockMutex.Unlock()
	}
}

func monitorStock(symbol string, scraper *StockScraper, hub *Hub, activeStocks *map[string]bool, mutex *sync.RWMutex) {
	log.Printf("Started monitoring %s", symbol)

	ticker := time.NewTicker(5 * time.Second) // Monitor each stock every 5 seconds
	defer ticker.Stop()

	for range ticker.C {
		mutex.RLock()
		isActive := (*activeStocks)[symbol]
		mutex.RUnlock()

		if !isActive {
			log.Printf("Stopped monitoring %s", symbol)
			mutex.Lock()
			delete(*activeStocks, symbol)
			mutex.Unlock()
			return
		}

		stockPrice, err := scraper.GetStockPrice(symbol)
		if err != nil {
			log.Printf("Error getting stock price for %s: %v", symbol, err)
			continue
		}

		log.Printf("[%s] %s Stock Price: $%s (Change: %s, %s%%)",
			stockPrice.Timestamp.Format("15:04:05"),
			symbol,
			stockPrice.Price,
			stockPrice.ChangeAbsolute,
			stockPrice.ChangePercent)

		// Broadcast to WebSocket clients
		hub.Broadcast(stockPrice)
	}
}
