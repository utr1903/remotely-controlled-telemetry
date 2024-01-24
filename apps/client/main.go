package main

import (
	"context"

	"github.com/utr1903/remotely-controlled-telemetry/apps/client/app"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/controller"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/otel"
)

func main() {

	ctx := context.Background()

	// Create metric provider
	mp := otel.NewMetricProvider(ctx)
	defer otel.ShutdownMetricProvider(ctx, mp)

	// Instantiate logger
	l := logger.New()

	// Run controller
	c := controller.New(l, "ws://localhost:8081/ws")
	go c.Run()

	// Run the application
	a := app.New(l)
	a.Run()
}
