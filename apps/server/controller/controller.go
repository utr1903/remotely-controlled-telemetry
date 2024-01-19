package controller

import (
	"fmt"
	"sync"
)

const WEB_SOCKET_PORT = "8081"

type Controller struct {
	wg        *sync.WaitGroup
	websocket *WebSocket
}

func New() *Controller {
	var wg sync.WaitGroup
	wg.Add(1)

	ws := newWebSocket(&wg, WEB_SOCKET_PORT)
	return &Controller{
		wg:        &wg,
		websocket: ws,
	}
}

func (c *Controller) Run() {

	go c.websocket.run()

	fmt.Println("Done")
	c.wg.Wait()
}
