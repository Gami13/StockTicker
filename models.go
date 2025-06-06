package main

import (
	"time"

	"github.com/gorilla/websocket"
)

// StockPrice represents the current stock price and change information
type StockPrice struct {
	Symbol         string    `json:"symbol"`
	Price          string    `json:"price"`
	ChangeAbsolute string    `json:"changeAbsolute"`
	ChangePercent  string    `json:"changePercent"`
	Timestamp      time.Time `json:"timestamp"`
}

// Config holds application configuration
type Config struct {
	Interval time.Duration
	Port     string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		Interval: 5 * time.Second,
		Port:     ":8081",
	}
}

// ClientMessage represents a message from the client
type ClientMessage struct {
	Type   string `json:"type"`
	Symbol string `json:"symbol"`
}

// Client represents a WebSocket client with their requested stock symbol
type Client struct {
	Conn   *websocket.Conn
	Symbol string
}
