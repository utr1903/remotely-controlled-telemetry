package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

const HTTP_SERVER_PORT = "8082"

type latency struct {
	duration time.Duration
	mutex    *sync.Mutex
}

type App struct {
	logger        *logger.Logger
	latency       *latency
	latencyMetric metric.Float64Histogram
}

func New(
	logger *logger.Logger,
) *App {

	// Create custom latency histogram
	latencyMetric, err := otel.GetMeterProvider().Meter("application").
		Float64Histogram(
			"application.latency",
			metric.WithUnit("ms"),
			metric.WithDescription("Measures the duration of each app handling"),
		)
	if err != nil {
		panic(err)
	}

	return &App{
		logger: logger,
		latency: &latency{
			duration: time.Second,
			mutex:    &sync.Mutex{},
		},
		latencyMetric: latencyMetric,
	}
}

func (a *App) Run() {

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	durationChannel := make(chan time.Duration)

	// Start HTTP server to change latency
	hs := newHttpServer(a.logger, durationChannel)
	go hs.serve()

	// Run the application
	go a.runApp()

	for {
		select {
		// Check for latency change
		case duration := <-durationChannel:
			a.logger.LogWithFields(
				logrus.InfoLevel,
				"Latency change request is received.",
				map[string]string{
					"component.name": "application",
				})

			// Set the new duration
			a.setLatencyDuration(duration)

			// Watch for sigterm
		case <-interrupt:
			a.logger.LogWithFields(
				logrus.ErrorLevel,
				"Interrupt received, shutting down the application...",
				map[string]string{
					"component.name": "application",
				})
			// cancel()
			return
		}
	}
}

func (a *App) runApp() {

	a.logger.LogWithFields(
		logrus.InfoLevel,
		"Running the application...",
		map[string]string{
			"component.name": "application",
		})

	for {

		// Start timer
		startTime := time.Now()

		// Sleep for given duration
		duration := a.getLatencyDuration()
		time.Sleep(duration)

		// Act as if the application is having some debug logs
		a.logger.LogWithFields(
			logrus.DebugLevel,
			"Some fancy debug log.",
			map[string]string{
				"component.name": "application",
			})

		attrs := make([]attribute.KeyValue, 0, 1)
		attrs = append(attrs, attribute.String("component.name", "application"))

		elapsedTime := float64(time.Since(startTime)) / float64(time.Millisecond)
		a.latencyMetric.Record(context.Background(), elapsedTime, metric.WithAttributes(attrs...))
	}
}

func (a *App) getLatencyDuration() time.Duration {
	a.latency.mutex.Lock()
	defer a.latency.mutex.Unlock()
	duration := a.latency.duration
	return duration
}

func (a *App) setLatencyDuration(
	duration time.Duration,
) {
	a.latency.mutex.Lock()
	defer a.latency.mutex.Unlock()
	a.latency.duration = duration
}
