package otel

import (
	"context"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

// Creates new meter provider
func NewMetricProvider(
	ctx context.Context,
) *sdkmetric.MeterProvider {

	// Create OTLP metric exporter
	exp, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(
			sdkmetric.NewPeriodicReader(exp, sdkmetric.WithInterval(
				10*time.Second,
			)),
		),
	)
	otel.SetMeterProvider(mp)
	return mp
}

// Shuts down meter provider
func ShutdownMetricProvider(
	ctx context.Context,
	mp *sdkmetric.MeterProvider,
) {
	// Do not make the application hang when it is shutdown.
	ctx, cancel := context.WithTimeout(ctx, time.Second*2)
	defer cancel()
	if err := mp.Shutdown(ctx); err != nil {
		panic(err)
	}
}
