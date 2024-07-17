package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
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

var connA *websocket.Conn

func handleConnections(w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatalf("Failed to upgrade to websocket: %v", err)
	}
	defer ws.Close()

	for {
		_, message, err := ws.ReadMessage()
		if err != nil {
			log.Printf("Error reading message: %v", err)
			break
		}
		log.Printf("Service B received: %s", message)
	}
}

func connectToServiceA() {
	var err error
	for {
		connA, _, err = websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
		if err == nil {
			break
		}
		log.Printf("Failed to connect to Service A: %v", err)
		time.Sleep(2 * time.Second)
	}
	log.Println("Connected to Service A")
}

func readConsoleAndSendMessages() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter message for Service A: ")
		message, _ := reader.ReadString('\n')

		err := connA.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Printf("Error writing message to Service A: %v", err)
			break
		}
	}
}

func main() {
	go connectToServiceA()

	http.HandleFunc("/ws", handleConnections)
	log.Println("Service B started on :8081")
	go func() {
		log.Fatal(http.ListenAndServe(":8081", nil))
	}()

	time.Sleep(3 * time.Second)
	go readConsoleAndSendMessages()

	select {}
}
