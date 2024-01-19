package controller

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type websocketClient struct {
	websocketServerUrl string
	wg                 *sync.WaitGroup
}

func newWebSocketClient(
	wg *sync.WaitGroup,
	websocketServerUrl string,
) *websocketClient {
	return &websocketClient{
		websocketServerUrl: websocketServerUrl,
		wg:                 wg,
	}
}

func (wc *websocketClient) run() {
	defer wc.wg.Done()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	fmt.Println("HERE")
	conn, _, err := websocket.DefaultDialer.Dial(wc.websocketServerUrl, nil)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("Error reading message:", err)
				return
			}
			fmt.Printf("Received message from server: %s\n", message)
		}
	}()

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			fmt.Println("Done1")
			return
		case <-ticker.C:
			// Do nothing, just wait for messages from the server
			fmt.Println("Ticket")
		case <-interrupt:
			fmt.Println("Interrupt received, closing connection.")
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				fmt.Println("Error sending close message:", err)
				return
			}
			select {
			case <-done:
				fmt.Println("Done2")
			case <-time.After(time.Second):
				fmt.Println("Ticker2")
			}
			return
		}
	}
}
