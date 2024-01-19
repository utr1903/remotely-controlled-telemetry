package controller

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type webSocketServer struct {
	controllerSwitch chan bool
	wg               *sync.WaitGroup
	port             string
	upgrader         *websocket.Upgrader
}

func newWebSocketServer(
	wg *sync.WaitGroup,
	cs chan bool,
	port string,
) *webSocketServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &webSocketServer{
		controllerSwitch: cs,
		wg:               wg,
		port:             port,
		upgrader:         &upgrader,
	}
}

func (ws *webSocketServer) run() {
	defer ws.wg.Done()

	http.HandleFunc("/ws", ws.handleConnections)

	fmt.Println("Web socket server is running on ", ws.port)
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
	fmt.Println("Web socket connection is established.")

	for turnOn := range ws.controllerSwitch {
		fmt.Println("Received: ", turnOn)
	}

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
