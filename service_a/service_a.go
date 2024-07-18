package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	connB *websocket.Conn
	mu    sync.Mutex
)

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade to websocket: %v", err)
		return
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Service A received: %s", message)
	}
}

func connectToServiceB() {
	var err error
	for {
		connB, _, err = websocket.DefaultDialer.Dial("ws://localhost:8081/ws", nil)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to Service B: %v", err)
		time.Sleep(2 * time.Second)
	}
	log.Println("Connected to Service B")
}

func readConsoleAndSendMessages() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message for Service B: ")
		message, _ := reader.ReadString('\n')

		mu.Lock()
		if connB != nil {
			err := connB.WriteMessage(websocket.TextMessage, []byte(message))
			if err != nil {
				log.Printf("Error writing message to Service B: %v", err)
				connB.Close()
				connB = nil
				go connectToServiceB()
			}
		}
		mu.Unlock()
	}
}

func main() {
	go connectToServiceB()

	http.HandleFunc("/ws", handleConnections)
	log.Println("Service A started on :8080")
	go func() {
		log.Fatal(http.ListenAndServe(":8080", nil))
	}()

	time.Sleep(3 * time.Second)
	go readConsoleAndSendMessages()

	select {}
}
