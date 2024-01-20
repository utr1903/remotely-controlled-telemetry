package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type webSocketServer struct {
	webSocketReadyChannel chan bool
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	port                  string
	upgrader              *websocket.Upgrader
}

func newWebSocketServer(
	wg *sync.WaitGroup,
	webSocketReadyChannel chan bool,
	controllerChannel chan bool,
	port string,
) *webSocketServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	return &webSocketServer{
		webSocketReadyChannel: webSocketReadyChannel,
		controllerChannel:     controllerChannel,
		wg:                    wg,
		port:                  port,
		upgrader:              &upgrader,
	}
}

func (ws *webSocketServer) run() {
	defer ws.wg.Done()

	http.HandleFunc("/ws", ws.handleConnections)

	fmt.Println("Web socket server is running on ", ws.port)
	err := http.ListenAndServe("localhost:"+ws.port, nil)
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

	// Reset socket connection to false
	conn.SetCloseHandler(
		func(code int, text string) error {
			fmt.Println("websocketserver: Web socket connection is lost.")
			ws.webSocketReadyChannel <- false
			return nil
		})

	fmt.Println("websocketserver: Web socket connection is established.")
	ws.webSocketReadyChannel <- true

	// Check for incoming messages in case the client disconnects
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				return
			}
		}
	}()

	for enable := range ws.controllerChannel {
		fmt.Println("websocketserver: Received: ", enable)
		var message []byte
		if enable {
			message = []byte("run")
		} else {
			message = []byte("stop")
		}
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
}
