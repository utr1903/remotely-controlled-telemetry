package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketServer struct {
	wg       *sync.WaitGroup
	port     string
	upgrader *websocket.Upgrader
}

func newWebSocketServer(
	wg *sync.WaitGroup,
	port string,
) *webSocketServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &webSocketServer{
		wg:       wg,
		port:     port,
		upgrader: &upgrader,
	}
}

func (ws *webSocketServer) run() {
	defer ws.wg.Done()

	http.HandleFunc("/ws", ws.handleConnections)
	fmt.Println("Server is running on ", ws.port)

	err := http.ListenAndServe(":"+ws.port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (ws *webSocketServer) handleConnections(
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
