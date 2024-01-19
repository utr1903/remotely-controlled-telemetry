package main

import "github.com/utr1903/remotely-controlled-telemetry/apps/server/controller"

func main() {

	// Run the controller
	c := controller.New()
	c.Run()
}
