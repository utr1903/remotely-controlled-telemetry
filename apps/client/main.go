package main

import (
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/otelcollector"
)

func main() {
	otelcol := otelcollector.New()
	otelcol.Run()

	// ctrl := controller.New("ws://localhost:8081/ws")
	// ctrl.Run()
}
