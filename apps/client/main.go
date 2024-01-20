package main

import (
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/controller"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

func main() {

	// Instantiate logger
	l := logger.New()

	// Run controller
	c := controller.New(l, "ws://localhost:8081/ws")
	c.Run()
}
