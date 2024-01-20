package controller

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/server/logger"
)

const HTTP_SERVER_PORT = "8080"
const WEB_SOCKET_PORT = "8081"

type Controller struct {
	logger                *logger.Logger
	webSocketReadyChannel chan bool
	controllerChannel     chan bool
	wg                    *sync.WaitGroup
	httpserver            *HttpServer
	websocketserver       *webSocketServer
}

func New(
	logger *logger.Logger,
) *Controller {
	webSocketReadyChannel := make(chan bool)
	controllerChannel := make(chan bool)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	hs := newHttpServer(logger, wg, webSocketReadyChannel, controllerChannel, HTTP_SERVER_PORT)
	ws := newWebSocketServer(logger, wg, webSocketReadyChannel, controllerChannel, WEB_SOCKET_PORT)

	return &Controller{
		logger:                logger,
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

	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Controller is started.",
		map[string]string{
			"component.name": "controller",
		})

	c.wg.Wait()

	close(c.webSocketReadyChannel)
	close(c.controllerChannel)
}
