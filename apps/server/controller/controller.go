package controller

import (
	"fmt"
	"sync"
)

const HTTP_SERVER_PORT = "8080"
const WEB_SOCKET_PORT = "8081"

type Controller struct {
	controllerSwitch chan bool
	wg               *sync.WaitGroup
	httpserver       *HttpServer
	websocketserver  *webSocketServer
}

func New() *Controller {
	cs := make(chan bool)
	wg := &sync.WaitGroup{}

	wg.Add(2)
	hs := newHttpServer(wg, cs, HTTP_SERVER_PORT)
	ws := newWebSocketServer(wg, cs, WEB_SOCKET_PORT)

	return &Controller{
		controllerSwitch: cs,
		wg:               wg,
		httpserver:       hs,
		websocketserver:  ws,
	}
}

func (c *Controller) Run() {

	go c.httpserver.run()
	go c.websocketserver.run()

	fmt.Println("Controller is started.")
	c.wg.Wait()

	close(c.controllerSwitch)
}
