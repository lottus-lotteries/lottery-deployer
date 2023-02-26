package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,

	CheckOrigin: func(r *http.Request) bool { return true },
}

func reader(conn *websocket.Conn) {
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		fmt.Println(string(p))

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

// Define our websocket endpoint
func serveWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.Host)

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	// Listen indefinitely through our websocket connection for new messages
	reader(ws)
}

// Setup endpoints
func setupRoutes() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Simple Server")
	})
	http.HandleFunc("/gen", serveWs)
}

func main() {
	err := generateNewLottery("FirstLottery", 1000)
	if err != nil {
		fmt.Println(err)
	}

	setupRoutes()
	http.ListenAndServe(":8080", nil)
}

func generateNewLottery(name string, tickets int) error {
	lotteryData := &ContractData{
		name,
		tickets,
	}
	lotteryEngine := NewEngine(lotteryData)

	err := lotteryEngine.GenerateWrapper()
	if err != nil {
		return fmt.Errorf("generating lottery: %w", err)
	}

	return nil
}
