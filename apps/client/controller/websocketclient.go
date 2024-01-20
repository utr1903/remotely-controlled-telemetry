package controller

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

type websocketClient struct {
	logger             *logger.Logger
	wg                 *sync.WaitGroup
	controllerChannel  chan bool
	websocketServerUrl string
}

func newWebSocketClient(
	logger *logger.Logger,
	wg *sync.WaitGroup,
	controllerChannel chan bool,
	websocketServerUrl string,
) *websocketClient {
	return &websocketClient{
		logger:             logger,
		wg:                 wg,
		controllerChannel:  controllerChannel,
		websocketServerUrl: websocketServerUrl,
	}
}

func (wc *websocketClient) run() {
	defer wc.wg.Done()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	wc.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting web socket client...",
		map[string]string{
			"component.name": "websocketclient",
		})

	conn, _, err := websocket.DefaultDialer.Dial(wc.websocketServerUrl, nil)
	if err != nil {
		wc.logger.LogWithFields(
			logrus.ErrorLevel,
			"Starting web socket client is failed: "+err.Error(),
			map[string]string{
				"component.name": "websocketclient",
			})
		return
	}
	defer conn.Close()

	done := make(chan struct{})

	go func() {
		defer close(done)
		defer close(wc.controllerChannel)

		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				wc.logger.LogWithFields(
					logrus.ErrorLevel,
					"Error occurred during reading message: "+err.Error(),
					map[string]string{
						"component.name": "websocketclient",
					})

				wc.controllerChannel <- false

				return
			}
			wc.logger.LogWithFields(
				logrus.InfoLevel,
				"Message is read: "+string(message),
				map[string]string{
					"component.name": "websocketclient",
				})

			msg := string(message)
			if msg == "run" {
				wc.controllerChannel <- true

			} else if msg == "stop" {
				wc.controllerChannel <- false
			}
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
			wc.logger.LogWithFields(
				logrus.ErrorLevel,
				"Interrupt received, closing connection...",
				map[string]string{
					"component.name": "websocketclient",
				})
			err := conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				wc.logger.LogWithFields(
					logrus.ErrorLevel,
					"Error occurred during sending close message: "+err.Error(),
					map[string]string{
						"component.name": "websocketclient",
					})
				return
			}
			select {
			case <-done:
				fmt.Println("Done2")
			case <-time.After(time.Second):
				wc.logger.LogWithFields(
					logrus.DebugLevel,
					"Health check successful.",
					map[string]string{
						"component.name": "websocketclient",
					})
			}
			return
		}
	}
}
