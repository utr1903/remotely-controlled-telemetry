package main

import (
	"github.com/utr1903/remotely-controlled-telemetry/apps/server/controller"
	"github.com/utr1903/remotely-controlled-telemetry/apps/server/logger"
)

func main() {

	// Instantiate logger
	l := logger.New()

	// Run the controller
	c := controller.New(l)
	c.Run()
}
