package otelcollector

import (
	"fmt"
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
	logger             *logger.Logger
	runnerSynchronizer *runnerSynchronizer
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
	}
}

func (c *Collector) Start() error {

	// Start collector
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Starting collector...",
		map[string]string{
			"component.name": "collector",
		})
	err := c.start()
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

	// // Wait for the process to finish or be interrupted
	// if err := cmd.Wait(); err != nil {
	// 	fmt.Println("Process finished with error:", err)
	// }

	// // Stop the process by sending SIGTERM
	// if err := Stop(*pid); err != nil {
	// 	fmt.Println("Error stopping process:", err)
	// }
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
	fmt.Println("Starting otel collector...")
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
		"OTel collector is started on process: "+strconv.FormatInt(int64(pid), 10),
		map[string]string{
			"component.name": "collector",
		})
	c.sync(true, &pid)

	return nil
}

func (c *Collector) Stop() error {
	// Get process ID
	pid := c.getPid()
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Stopping OTel collector on process: "+strconv.FormatInt(int64(pid), 10),
		map[string]string{
			"component.name": "collector",
		})

	// Find the process by its ID
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Finding process: "+strconv.FormatInt(int64(pid), 10),
		map[string]string{
			"component.name": "collector",
		})
	process, err := os.FindProcess(*c.runnerSynchronizer.pid)
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Process "+strconv.FormatInt(int64(pid), 10)+" is not found: "+err.Error(),
			map[string]string{
				"component.name": "collector",
			})
		return err
	}

	// Send SIGTERM signal to the process
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Process "+strconv.FormatInt(int64(pid), 10)+" is found. Stopping...",
		map[string]string{
			"component.name": "collector",
		})
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		c.logger.LogWithFields(
			logrus.ErrorLevel,
			"Stopping process "+strconv.FormatInt(int64(pid), 10)+" failed: "+err.Error(),
			map[string]string{
				"component.name": "collector",
			})
		return err
	}
	c.logger.LogWithFields(
		logrus.InfoLevel,
		"Stopping process "+strconv.FormatInt(int64(pid), 10)+" succeeded: "+err.Error(),
		map[string]string{
			"component.name": "collector",
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
