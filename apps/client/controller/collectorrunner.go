package controller

import (
	"sync"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/otelcollector"
)

type collectorRunner struct {
	logger            *logger.Logger
	wg                *sync.WaitGroup
	controllerChannel chan bool
	otelcol           *otelcollector.Collector
}

func newCollectorRunner(
	logger *logger.Logger,
	wg *sync.WaitGroup,
	controllerChannel chan bool,
) *collectorRunner {
	otelcol := otelcollector.New(logger)

	return &collectorRunner{
		logger:            logger,
		wg:                wg,
		controllerChannel: controllerChannel,
		otelcol:           otelcol,
	}
}

func (cr *collectorRunner) run() {
	defer cr.wg.Done()

	cr.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting controller runner...",
		map[string]string{
			"component.name": "controllerrunner",
		})
	err := cr.otelcol.Start()
	if err != nil {
		cr.logger.LogWithFields(
			logrus.ErrorLevel,
			"Starting controller runner is failed.",
			map[string]string{
				"component.name": "controllerrunner",
				"error.message":  err.Error(),
			})
		return
	}

	cr.logger.LogWithFields(
		logrus.InfoLevel,
		"Controller runner is started. Listening controller channel...",
		map[string]string{
			"component.name": "controllerrunner",
		})

	for run := range cr.controllerChannel {
		if run {
			cr.otelcol.Start()
		} else {
			cr.otelcol.Stop()
		}
	}
}
