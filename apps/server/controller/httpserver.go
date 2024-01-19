package controller

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type HttpServer struct {
	wg       *sync.WaitGroup
	port     string
	upgrader *websocket.Upgrader
}

func newHttpServer(
	wg *sync.WaitGroup,
	port string,
) *HttpServer {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	return &HttpServer{
		wg:       wg,
		port:     port,
		upgrader: &upgrader,
	}
}

func (hs *HttpServer) run() {
	defer hs.wg.Done()

	http.Handle("/control", http.HandlerFunc(hs.handleTelemetryCollection))
	http.ListenAndServe(":"+hs.port, nil)
}

func (hs *HttpServer) handleTelemetryCollection(
	w http.ResponseWriter,
	r *http.Request,
) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
