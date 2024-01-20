package controller

import (
	"fmt"
	"sync"
)

type Controller struct {
	controllerChannel chan bool
	wg                *sync.WaitGroup
	webSocketClient   *websocketClient
	collectorRunner   *collectorRunner
}

func New(
	webSocketUrl string,
) *Controller {

	controllerChannel := make(chan bool)

	wg := &sync.WaitGroup{}

	wg.Add(2)
	cr := newCollectorRunner(wg, controllerChannel)
	wc := newWebSocketClient(wg, controllerChannel, webSocketUrl)

	return &Controller{
		controllerChannel: controllerChannel,
		wg:                wg,
		webSocketClient:   wc,
		collectorRunner:   cr,
	}
}

func (c *Controller) Run() {

	go c.collectorRunner.run()
	go c.webSocketClient.run()

	fmt.Println("Controller is started.")
	c.wg.Wait()
}
