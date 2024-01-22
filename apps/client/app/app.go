package app

import (
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

const HTTP_SERVER_PORT = "8082"

type latency struct {
	duration time.Duration
	mutex    *sync.Mutex
}

type App struct {
	logger  *logger.Logger
	latency *latency
}

func New(
	logger *logger.Logger,
) *App {
	return &App{
		logger: logger,
		latency: &latency{
			duration: time.Second,
			mutex:    &sync.Mutex{},
		},
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
				logrus.DebugLevel,
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
