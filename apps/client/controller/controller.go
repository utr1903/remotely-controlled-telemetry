package controller

import (
	"fmt"
	"sync"
)

type Controller struct {
	wg              *sync.WaitGroup
	webSocketClient *websocketClient
}

func New(
	webSocketUrl string,
) *Controller {
	wg := &sync.WaitGroup{}

	wg.Add(1)
	wc := newWebSocketClient(wg, webSocketUrl)

	return &Controller{
		wg:              wg,
		webSocketClient: wc,
	}
}

func (c *Controller) Run() {

	go c.webSocketClient.run()

	fmt.Println("Controller is started.")
	c.wg.Wait()
}
