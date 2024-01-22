package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/server/logger"
)

type webSocketSynchronizer struct {
	isReady bool
	mutex   sync.Mutex
}

type HttpServer struct {
	logger                *logger.Logger
	webSocketReadyChannel chan bool
	webSocketSynchronizer *webSocketSynchronizer
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	port                  string
}

func newHttpServer(
	logger *logger.Logger,
	wg *sync.WaitGroup,
	webSocketReadyChannel chan bool,
	controllerChannel chan bool,
	port string,
) *HttpServer {

	webSocketSynchronizer := &webSocketSynchronizer{
		isReady: false,
		mutex:   sync.Mutex{},
	}

	return &HttpServer{
		logger:                logger,
		webSocketReadyChannel: webSocketReadyChannel,
		webSocketSynchronizer: webSocketSynchronizer,
		controllerChannel:     controllerChannel,
		wg:                    wg,
		port:                  port,
	}
}

func (hs *HttpServer) run() {
	defer hs.wg.Done()

	go hs.synchronizeWebSocketConnection()

	http.Handle("/control", http.HandlerFunc(hs.handleTelemetryCollection))

	hs.logger.LogWithFields(
		logrus.InfoLevel,
		"HTTP server is running on localhost:"+hs.port,
		map[string]string{
			"component.name": "httpserver",
		})
	err := http.ListenAndServe("localhost:"+hs.port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (hs *HttpServer) handleTelemetryCollection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if !hs.isWebSocketReady() {
		msg := "Web socket connection is not yet established!"
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "httpserver",
			})
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(msg))
		return
	}

	switch r.Method {
	case http.MethodPost:
		hs.controllerChannel <- true

		msg := "Signal is sent to the client to run the collector."
		hs.logger.LogWithFields(
			logrus.InfoLevel,
			msg,
			map[string]string{
				"component.name": "httpserver",
			})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))

	case http.MethodDelete:
		hs.controllerChannel <- false

		msg := "Signal is sent to the client to stop the collector."
		hs.logger.LogWithFields(
			logrus.InfoLevel,
			msg,
			map[string]string{
				"component.name": "httpserver",
			})
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(msg))

	default:
		msg := "Request is not valid!"
		hs.logger.LogWithFields(
			logrus.ErrorLevel,
			msg,
			map[string]string{
				"component.name": "httpserver",
			})
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(msg))
		return
	}
}

func (hs *HttpServer) synchronizeWebSocketConnection() {
	for isReady := range hs.webSocketReadyChannel {
		if isReady {
			hs.logger.LogWithFields(
				logrus.InfoLevel,
				"Web socket connection is established.",
				map[string]string{
					"component.name": "httpserver",
				})

		} else {
			hs.logger.LogWithFields(
				logrus.ErrorLevel,
				"Web socket connection is lost.",
				map[string]string{
					"component.name": "httpserver",
				})
			return
		}

		hs.webSocketSynchronizer.mutex.Lock()
		hs.webSocketSynchronizer.isReady = isReady
		hs.webSocketSynchronizer.mutex.Unlock()
	}
}

func (hs *HttpServer) isWebSocketReady() bool {
	hs.webSocketSynchronizer.mutex.Lock()
	defer hs.webSocketSynchronizer.mutex.Unlock()
	return hs.webSocketSynchronizer.isReady
}
