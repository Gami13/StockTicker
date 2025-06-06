package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// StockScraper handles stock price scraping from Bing
type StockScraper struct {
	allocCtx context.Context
}

// NewStockScraper creates a new stock scraper with Chrome allocator context
func NewStockScraper() (*StockScraper, error) {
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", false),
		chromedp.Flag("disable-dev-shm-usage", true),
		chromedp.Flag("no-sandbox", true),
		chromedp.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"),
		chromedp.WindowSize(1920, 1080),
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.Flag("disable-features", "VizDisplayCompositor"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	return &StockScraper{allocCtx: allocCtx}, nil
}

// GetStockPrice scrapes the current stock price for a specific symbol
func (s *StockScraper) GetStockPrice(symbol string) (StockPrice, error) {
	var stockPrice StockPrice
	var changeString string

	// Create a new context for this request to avoid conflicts
	ctx, cancel := chromedp.NewContext(s.allocCtx)
	defer cancel()

	// Set a timeout for the scraping operation
	ctx, timeoutCancel := context.WithTimeout(ctx, 15*time.Second)
	defer timeoutCancel()

	// Navigate to the specific stock's page
	url := fmt.Sprintf("https://www.bing.com/search?q=%s+stock+price", symbol)
	err := chromedp.Run(ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("body"),
		chromedp.WaitReady(".b_focusTextMedium", chromedp.ByQuery),
		chromedp.Sleep(2*time.Second),
	)
	if err != nil {
		return stockPrice, fmt.Errorf("error navigating to %s stock page: %w", symbol, err)
	}

	// Get the stock price
	err = chromedp.Run(ctx,
		chromedp.Text(".b_focusTextMedium", &stockPrice.Price, chromedp.ByQuery),
	)
	if err != nil || stockPrice.Price == "" {
		return stockPrice, fmt.Errorf("error getting stock price for %s: %w", symbol, err)
	}

	// Get the change information
	err = chromedp.Run(ctx,
		chromedp.Text(".fin_change", &changeString, chromedp.ByQuery),
	)
	if err != nil {
		log.Printf("Warning: Could not get change information for %s: %v", symbol, err)
	}

	// Parse change information
	if changeString != "" {
		s.parseChangeString(changeString, &stockPrice)
	}

	stockPrice.Symbol = symbol
	stockPrice.Timestamp = time.Now()

	return stockPrice, nil
}

// parseChangeString parses the change string from Bing (e.g., "▼ -47,35 (-14,26%) today")
func (s *StockScraper) parseChangeString(changeString string, stockPrice *StockPrice) {
	parts := strings.Fields(changeString)
	if len(parts) >= 3 {
		// parts[1] should be the absolute change
		stockPrice.ChangeAbsolute = parts[1]

		// parts[2] should be the percentage in parentheses
		if strings.HasPrefix(parts[2], "(") && strings.HasSuffix(parts[2], "%)") {
			stockPrice.ChangePercent = strings.TrimSuffix(strings.TrimPrefix(parts[2], "("), "%)")
		}
	}
}

// Close closes the scraper context
func (s *StockScraper) Close() {
	// Note: In a real implementation, you'd want to properly manage the context lifecycle
}
