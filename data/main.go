package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"
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

// Serve stock data over HTTP
func stockHandler(w http.ResponseWriter, r *http.Request) {
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		http.Error(w, "Symbol is required", http.StatusBadRequest)
		return
	}

	mu.Lock()
	data, exists := stockStore[symbol]
	mu.Unlock()

	if !exists {
		http.Error(w, "Stock data not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
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
	http.HandleFunc("/stock", stockHandler)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	wg.Wait()
}
