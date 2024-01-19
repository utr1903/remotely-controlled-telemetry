package controller

import (
	"fmt"
	"sync"
)

const HTTP_SERVER_PORT = "8080"
const WEB_SOCKET_PORT = "8081"

type Controller struct {
	wg              *sync.WaitGroup
	httpserver      *HttpServer
	websocketserver *webSocketServer
}

func New() *Controller {
	var wg sync.WaitGroup

	wg.Add(1)
	hs := newHttpServer(&wg, HTTP_SERVER_PORT)

	wg.Add(1)
	ws := newWebSocketServer(&wg, WEB_SOCKET_PORT)

	return &Controller{
		wg:              &wg,
		httpserver:      hs,
		websocketserver: ws,
	}
}

func (c *Controller) Run() {

	go c.httpserver.run()
	go c.websocketserver.run()

	fmt.Println("Done")
	c.wg.Wait()
}
