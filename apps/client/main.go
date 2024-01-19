package main

import "github.com/utr1903/remotely-controlled-telemetry/apps/client/controller"

func main() {
	c := controller.New("ws://localhost:8081/ws")
	c.Run()
}
