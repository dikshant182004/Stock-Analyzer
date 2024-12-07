package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// StockData represents the stock information
type StockData struct {
	Symbol      string  `json:"symbol"`
	Price       float64 `json:"price"`
	High        float64 `json:"high"`
	Low         float64 `json:"low"`
	LastUpdated string  `json:"lastUpdated"`
	Error       string  `json:"error,omitempty"`
}

// Global storage for stock data
var stockStore = make(map[string][]StockData)
var mu sync.Mutex

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // all origins
	},
}

// using websocket instead of http
func stockWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Printf("WebSocket upgrade error:")
	}
	defer conn.Close()

	// Get stock symbol from query
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		log.Println("No stock symbol provided")
		// Send a message to the client if the symbol is missing
		conn.WriteMessage(websocket.TextMessage, []byte("Error: No stock symbol provided"))
		return
	}

	// Stream stock data to WebSocket client
	for {
		mu.Lock()
		data, exists := stockStore[symbol]
		mu.Unlock()

		if !exists {
			data = append(data, StockData{
				Symbol:      symbol,
				Price:       0.0,
				LastUpdated: time.Now().Format(time.RFC3339),
			})
		}

		if err := conn.WriteJSON(data); err != nil {
			log.Printf("Error writing WebSocket message: %v", err)
			return
		}

		time.Sleep(1 * time.Second) // Adjust as needed for real-time updates
	}

}

func fetchStock(symbol string) {
	for {
		cmd := exec.Command("C:/Users/Dell/AppData/Local/Programs/Python/Python311/python.exe", "J:/personel/StockAnalyzer/data/real-time/fetch_stock.py", symbol)
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("Error fetching stock %s: %v, Output: %s", symbol, err, string(output))
			time.Sleep(10 * time.Second)
			continue
		}

		var data StockData // creating a variable of above struct
		if !json.Valid(output) {
			log.Printf("Invalid JSON for stock %s: %s", symbol, string(output))
			time.Sleep(10 * time.Second)
			continue
		}

		// just encoding the json data into go struct var
		if err := json.Unmarshal(output, &data); err != nil {
			log.Printf("Error decoding JSON for %s: %v", symbol, err)
			time.Sleep(10 * time.Second)
			continue
		}

		if data.Error != "" {
			log.Printf("Error from Python for %s: %s", symbol, data.Error)
			time.Sleep(10 * time.Second)
			continue
		}

		mu.Lock()
		// just storing the previously fetched data for visualization stuff
		stockStore[symbol] = append(stockStore[symbol], data)
		if len(stockStore[symbol]) > 100 {
			stockStore[symbol] = stockStore[symbol][1:]
		}
		mu.Unlock()

		log.Printf("Updated stock data for %s: %+v", symbol, data)
		time.Sleep(10 * time.Second) // Adjust as needed
	}
}

func main() {

	symbols := []string{"500112", "500325", "532540"}

	// making go routine for every obtained symbol

	wg := &sync.WaitGroup{}
	for _, symbol := range symbols {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			fetchStock(symbol)
		}(symbol)
	}

	// Start HTTP server to serve data
	http.HandleFunc("/ws", stockWebSocketHandler)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	wg.Wait()
}
