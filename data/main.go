package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os/exec"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
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

// a global variable representing the redis client to connect and interact with redis server

var (
	redisClient *redis.Client
	upgrader    = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow all origins
		},
	}
)

// WebSocket handler to subscribe to stock updates
func stockWebSocketHandler(w http.ResponseWriter, r *http.Request) {
	// upgrading http to websocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}
	defer conn.Close()

	// Get stock symbol from query
	symbol := r.URL.Query().Get("symbol")
	if symbol == "" {
		log.Println("No stock symbol provided")
		conn.WriteMessage(websocket.TextMessage, []byte("Error: No stock symbol provided"))
		return
	}

	// Subscribe to the Redis channel for the stock symbol to receive messages from it
	pubsub := redisClient.Subscribe(context.Background(), symbol) // publisher / Subscriber
	defer pubsub.Close()

	// Start receiving messages from Redis and send them to the WebSocket client
	for {
		msg, err := pubsub.ReceiveMessage(context.Background())
		if err != nil {
			log.Printf("Error receiving message from Redis: %v", err)
			return
		}

		// Send the message to the WebSocket client
		if err := conn.WriteMessage(websocket.TextMessage, []byte(msg.Payload)); err != nil {
			log.Printf("Error writing WebSocket message: %v", err)
			return
		}
		log.Printf("Received from Redis: channel=%s, message=%s", msg.Channel, msg.Payload)
	}
}

// Fetch stock data and publish updates to Redis
func fetchStock(symbol string) {
	for {
		cmd := exec.Command("C:/Users/Dell/AppData/Local/Programs/Python/Python311/python.exe", "J:/personel/StockAnalyzer/data/real-time/fetch_stock.py", symbol)
		output, err := cmd.CombinedOutput()

		if err != nil {
			log.Printf("Error fetching stock %s: %v, Output: %s", symbol, err, string(output))
			time.Sleep(10 * time.Second)
			continue
		}

		var data StockData
		if !json.Valid(output) {
			log.Printf("Invalid JSON for stock %s: %s", symbol, string(output))
			time.Sleep(10 * time.Second)
			continue
		}

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

		// Publish stock data to Redis channel
		dataJSON, _ := json.Marshal(data)
		// publish -any client or process sends messages to a channel(virtual room {where messages
		// are sent and received}) with the named topic :-symbol

		// ~ here the symbol name is created as a channel and messages are sent in it
		err = redisClient.Publish(context.Background(), symbol, dataJSON).Err()
		if err != nil {
			log.Printf("Error publishing to Redis channel %s: %v", symbol, err)
		}

		// log.Printf("Updated stock data for %s: %+v", symbol, data)
		time.Sleep(10 * time.Second) // Adjust as needed
	}
}

func main() {
	// Initialize Redis client
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0, // Use default DB
	})

	// Test Redis connection
	if _, err := redisClient.Ping(context.Background()).Result(); err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	}
	log.Println("Connected to Redis")

	// List of stock symbols to monitor
	symbols := []string{"500112", "500325", "532540"}

	// Fetch stock data in separate goroutines
	wg := &sync.WaitGroup{}
	for _, symbol := range symbols {
		wg.Add(1)
		go func(s string) {
			defer wg.Done()
			fetchStock(s)
		}(symbol)
	}

	// Start WebSocket server
	http.HandleFunc("/ws", stockWebSocketHandler)
	fmt.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))

	wg.Wait()
}
