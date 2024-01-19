package controller

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HttpServer struct {
	controllerSwitch chan bool
	wg               *sync.WaitGroup
	port             string
	upgrader         *websocket.Upgrader
}

func newHttpServer(
	wg *sync.WaitGroup,
	cs chan bool,
	port string,
) *HttpServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &HttpServer{
		controllerSwitch: cs,
		wg:               wg,
		port:             port,
		upgrader:         &upgrader,
	}
}

func (hs *HttpServer) run() {
	defer hs.wg.Done()

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
	if r.Method == http.MethodPost {
		fmt.Println("Turning on...")
		hs.controllerSwitch <- true
		fmt.Println("Turned on.")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Turned on"))
	} else if r.Method == http.MethodDelete {
		fmt.Println("Turning off...")
		hs.controllerSwitch <- false
		fmt.Println("Turned off.")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Turned off"))
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not valid!"))
	}
}
