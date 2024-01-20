package controller

import (
	"sync"

	"github.com/utr1903/remotely-controlled-telemetry/apps/client/otelcollector"
)

type collectorRunner struct {
	wg                *sync.WaitGroup
	controllerChannel chan bool
	otelcol           *otelcollector.Collector
}

func newCollectorRunner(
	wg *sync.WaitGroup,
	controllerChannel chan bool,
) *collectorRunner {
	otelcol := otelcollector.New()

	return &collectorRunner{
		wg:                wg,
		controllerChannel: controllerChannel,
		otelcol:           otelcol,
	}
}

func (cr *collectorRunner) run() {
	defer cr.wg.Done()

	err := cr.otelcol.Start()
	if err != nil {
		return
	}

	for run := range cr.controllerChannel {
		if run {
			cr.otelcol.Start()
		} else {
			cr.otelcol.Stop()
		}
	}
}
