package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type WebSocket struct {
	wg       *sync.WaitGroup
	port     string
	upgrader *websocket.Upgrader
}

func newWebSocket(
	wg *sync.WaitGroup,
	port string,
) *WebSocket {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &WebSocket{
		wg:       wg,
		port:     port,
		upgrader: &upgrader,
	}
}

func (ws *WebSocket) run() {
	defer ws.wg.Done()

	http.HandleFunc("/ws", ws.handleConnections)
	fmt.Println("Server is running on ", ws.port)

	err := http.ListenAndServe(":"+ws.port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (ws *WebSocket) handleConnections(
	w http.ResponseWriter,
	r *http.Request,
) {
	conn, err := ws.upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()

	for {
		// Send a message every 5 seconds
		time.Sleep(1 * time.Second)
		message := []byte("Hello from server!")
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
