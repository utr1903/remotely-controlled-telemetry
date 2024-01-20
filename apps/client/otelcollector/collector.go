package otelcollector

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"sync"
	"syscall"

	"github.com/sirupsen/logrus"
	"github.com/utr1903/remotely-controlled-telemetry/apps/client/logger"
)

type runnerSynchronizer struct {
	isRunning bool
	pid       *int
	mutex     *sync.Mutex
}

type Collector struct {
	logger                       *logger.Logger
	runnerSynchronizer           *runnerSynchronizer
	otelCollectorConfigGenerator *otelCollectorConfigGenerator
}

func New(
	logger *logger.Logger,
) *Collector {
	return &Collector{
		logger: logger,
		runnerSynchronizer: &runnerSynchronizer{
			isRunning: false,
			pid:       nil,
			mutex:     &sync.Mutex{},
		},
		otelCollectorConfigGenerator: newOtelCollectorConfigGenerator(
			logger,
		),
	}
}

func (c *Collector) Start() error {

	// Generate OTel collector config file
	err := c.otelCollectorConfigGenerator.generate()
	if err != nil {
		return err
	}

	// Start collector
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting collector...",
		map[string]string{
			"component.name": "collector",
		})
	err = c.start()
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Starting collector is failed: "+err.Error(),
			map[string]string{
				"component.name": "collector",
			})
		return err
	}

	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting collector is succeeded.",
		map[string]string{
			"component.name": "collector",
		})
	return nil
}

func (c *Collector) start() error {
	currentDir, err := os.Getwd()
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Current directory is not retrieved: "+err.Error(),
			map[string]string{
				"component.name": "collector",
			})
		return err
	}

	// Build the full path to the "app" executable in the current directory
	appPath := filepath.Join(currentDir, "/bin/otelcol-contrib")

	// Create a new Cmd struct for the "app" executable with the argument
	cmd := exec.Command(appPath, "--config=./bin/otel-config.yaml")

	// Start the process
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting OTel collector...",
		map[string]string{
			"component.name": "collector",
		})
	if err := cmd.Start(); err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Starting OTel collector is failed:"+err.Error(),
			map[string]string{
				"component.name": "collector",
			})
		return err
	}

	// Get the process ID
	pid := cmd.Process.Pid
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"OTel collector is started.",
		map[string]string{
			"component.name":     "collector",
			"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
		})
	c.sync(true, &pid)

	return nil
}

func (c *Collector) Stop() error {
	// Get process ID
	pid := c.getPid()
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Stopping OTel collector...",
		map[string]string{
			"component.name":     "collector",
			"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
		})

	// Find the process by its ID
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Finding process...",
		map[string]string{
			"component.name":     "collector",
			"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
		})
	process, err := os.FindProcess(*c.runnerSynchronizer.pid)
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Process is not found: "+err.Error(),
			map[string]string{
				"component.name":     "collector",
				"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
			})
		return err
	}

	// Send SIGTERM signal to the process
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Process is found. Stopping...",
		map[string]string{
			"component.name":     "collector",
			"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
		})
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Stopping process failed: "+err.Error(),
			map[string]string{
				"component.name":     "collector",
				"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
			})
		return err
	}
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Stopping process succeeded.",
		map[string]string{
			"component.name":     "collector",
			"otelcol.process.id": strconv.FormatInt(int64(pid), 10),
		})

	c.sync(false, nil)

	return nil
}

func (c *Collector) sync(
	isRunning bool,
	pid *int,
) {
	c.runnerSynchronizer.mutex.Lock()
	defer c.runnerSynchronizer.mutex.Unlock()
	c.runnerSynchronizer.isRunning = isRunning
	c.runnerSynchronizer.pid = pid
}

func (c *Collector) getPid() int {
	c.runnerSynchronizer.mutex.Lock()
	defer c.runnerSynchronizer.mutex.Unlock()
	pid := c.runnerSynchronizer.pid
	return *pid
}
