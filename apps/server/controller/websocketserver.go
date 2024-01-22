package controller

import (
	"fmt"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/server/logger"
)

type webSocketServer struct {
	logger                *logger.Logger
	webSocketReadyChannel chan bool
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	port                  string
	upgrader              *websocket.Upgrader
}

func newWebSocketServer(
	logger *logger.Logger,
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
		logger:                logger,
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

	ws.logger.LogWithFields(
		logrus.InfoLevel,
		"Web socket server is running on localhost:"+ws.port,
		map[string]string{
			"component.name": "websocketserver",
		})

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
			ws.logger.LogWithFields(
				logrus.ErrorLevel,
				"Web socket connection is lost.",
				map[string]string{
					"component.name": "websocketserver",
				})
			ws.webSocketReadyChannel <- false
			return nil
		})

	ws.logger.LogWithFields(
		logrus.InfoLevel,
		"Web socket connection is established.",
		map[string]string{
			"component.name": "websocketserver",
		})

	ws.webSocketReadyChannel <- true

	// Check for incoming messages in case the client disconnects
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				ws.logger.LogWithFields(
					logrus.ErrorLevel,
					"Error occurred during reading message from the client.",
					map[string]string{
						"component.name": "websocketserver",
						"error.message":  err.Error(),
					})
				return
			}
		}
	}()

	for otelColSwitch := range ws.controllerChannel {
		ws.logger.LogWithFields(
			logrus.InfoLevel,
			"Controller channel input received.",
			map[string]string{
				"component.name": "websocketserver",
				"otelcol.enable": strconv.FormatBool(otelColSwitch),
			})

		var message []byte
		if otelColSwitch {
			message = []byte("start")
		} else {
			message = []byte("stop")
		}
		err := conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			ws.logger.LogWithFields(
				logrus.ErrorLevel,
				"Error occurred during writing message to web socket.",
				map[string]string{
					"component.name": "websocketserver",
					"error.message":  err.Error(),
				})
		}
	}
}
