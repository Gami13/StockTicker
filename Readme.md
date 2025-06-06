# Stock Ticker for OBS

A real-time stock price ticker application designed specifically for embedding in OBS Studio. This application provides live stock price updates with a clean interface perfect for streaming overlays.

## Features

- **Real-time stock price updates** - Updates every 5 seconds
- **WebSocket-based live streaming** - Instant price updates without page refresh
- **OBS-ready design** - Clean, minimal interface optimized for streaming overlays
- **Multiple stock support** - Switch between different stocks via URL parameters
- **Price change indicators** - Visual indicators for positive/negative price changes
- **Auto-reconnection** - Automatically reconnects if connection is lost

## Screenshots

The ticker displays:
- Stock symbol
- Current price
- Absolute price change
- Percentage change
- Color-coded indicators (green for gains, red for losses)

## Installation

### Prerequisites

- Google Chrome or Chromium browser (for web scraping)

## Usage

### Basic Usage

1. Start the application by running the .exe from releases
  

2. Open your browser and navigate to:
   ```
   http://localhost:8081?stock=AAPL
   ```
   Replace `AAPL` with any valid stock symbol (e.g., `MSFT`, `GOOGL`, `TSLA`, etc.)

### OBS Studio Integration

1. **Add Browser Source:**
   - In OBS Studio, add a new "Browser Source"
   - Set the URL to: `http://localhost:8081?stock=SYMBOL`
   - Replace `SYMBOL` with your desired stock (e.g., `AAPL`)

2. **Recommended Settings:**
   - Width: 500px
   - Height: 70px
   - Check "Shutdown source when not visible" for better performance
   - Check "Refresh browser when scene becomes active" for reliability
   - The size is big for improved readability, you should scale it down in OBS to fit your layout







## How It Works

1. **Web Scraping**: Uses ChromeDP to scrape stock prices from Bing search results
2. **WebSocket Communication**: Real-time updates pushed to connected clients
3. **Multi-Stock Support**: Dynamically starts/stops monitoring based on client requests
4. **Concurrent Processing**: Each stock is monitored in a separate goroutine

## Development

### Default Settings

- **Server Port**: 8081
- **Update Interval**: 5 seconds per stock
- **WebSocket Endpoint**: `/ws`
- **Static Files**: Served from `/`
  
You can modify the appearance by editing `index.html`:







### Dependencies

- **gorilla/websocket**: WebSocket communication
- **chromedp/chromedp**: Headless Chrome automation for web scraping








**Note**: This application scrapes stock data from public sources. Please be respectful of rate limits and terms of service. For production use with high frequency updates, consider using a proper financial data API.