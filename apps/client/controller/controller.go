package controller

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

type Controller struct {
	logger            *logger.Logger
	controllerChannel chan bool
	wg                *sync.WaitGroup
	webSocketClient   *websocketClient
	collectorRunner   *collectorRunner
}

func New(
	logger *logger.Logger,
	webSocketUrl string,
) *Controller {

	controllerChannel := make(chan bool)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	cr := newCollectorRunner(logger, wg, controllerChannel)
	wc := newWebSocketClient(logger, wg, controllerChannel, webSocketUrl)

	return &Controller{
		logger:            logger,
		controllerChannel: controllerChannel,
		wg:                wg,
		webSocketClient:   wc,
		collectorRunner:   cr,
	}
}

func (c *Controller) Run() {

	go c.collectorRunner.run()
	go c.webSocketClient.run()

	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Controller is started.",
		map[string]string{
			"component.name": "controller",
		})
	c.wg.Wait()
}
