package controller

import (
	"fmt"
	"sync"
)

const HTTP_SERVER_PORT = "8080"
const WEB_SOCKET_PORT = "8081"

type Controller struct {
	webSocketReadyChannel chan bool
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	httpserver            *HttpServer
	websocketserver       *webSocketServer
}

func New() *Controller {
	webSocketReadyChannel := make(chan bool)
	controllerChannel := make(chan bool)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	hs := newHttpServer(wg, webSocketReadyChannel, controllerChannel, HTTP_SERVER_PORT)
	ws := newWebSocketServer(wg, webSocketReadyChannel, controllerChannel, WEB_SOCKET_PORT)

	return &Controller{
		webSocketReadyChannel: webSocketReadyChannel,
		controllerChannel:     controllerChannel,
		wg:                    wg,
		httpserver:            hs,
		websocketserver:       ws,
	}
}

func (c *Controller) Run() {

	go c.httpserver.run()
	go c.websocketserver.run()

	fmt.Println("Controller is started.")
	c.wg.Wait()

	close(c.webSocketReadyChannel)
	close(c.controllerChannel)
}
