package main

import "time"

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
	StockSymbol string
	Interval    time.Duration
	Port        string
}

// NewConfig creates a new configuration with default values
func NewConfig() *Config {
	return &Config{
		StockSymbol: "AAPL",
		Interval:    5 * time.Second,
		Port:        ":8080",
	}
}
