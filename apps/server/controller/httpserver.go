package controller

import (
	"fmt"
	"net/http"
	"sync"
)

type webSocketSynchronizer struct {
	isReady bool
	mutex   sync.Mutex
}

type HttpServer struct {
	webSocketReadyChannel chan bool
	webSocketSynchronizer *webSocketSynchronizer
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	port                  string
}

func newHttpServer(
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

	fmt.Println("HTTP server is running on ", hs.port)
	err := http.ListenAndServe(":"+hs.port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func (hs *HttpServer) handleTelemetryCollection(
	w http.ResponseWriter,
	r *http.Request,
) {

	if !hs.isWebSocketReady() {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Web socket connection is not yet established!"))
		return
	}

	switch r.Method {
	case http.MethodPost:
		hs.controllerChannel <- true

	case http.MethodDelete:
		hs.controllerChannel <- false

	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Request is not valid!"))
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Turned on"))
}

func (hs *HttpServer) synchronizeWebSocketConnection() {
	for isReady := range hs.webSocketReadyChannel {
		if isReady {
			fmt.Println("httpserver: Web socket connection is established.")
		} else {
			fmt.Println("httpserver: Web socket connection is lost.")
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
