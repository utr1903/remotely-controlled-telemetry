package otelcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"
)

type runnerSynchronizer struct {
	isRunning bool
	pid       *int
	mutex     *sync.Mutex
}

type Collector struct {
	runnerSynchronizer *runnerSynchronizer
}

func New() *Collector {
	return &Collector{
		runnerSynchronizer: &runnerSynchronizer{
			isRunning: false,
			pid:       nil,
			mutex:     &sync.Mutex{},
		},
	}
}

func (c *Collector) Start() error {

	// Start collector
	err := c.start()
	if err != nil {
		return err
	}
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
		fmt.Println("Error getting current directory:", err)
		return err
	}

	// Build the full path to the "app" executable in the current directory
	appPath := filepath.Join(currentDir, "/bin/otelcol-contrib")

	// Create a new Cmd struct for the "app" executable with the argument
	cmd := exec.Command(appPath, "--config=./bin/otel-config.yaml")

	// Start the process
	fmt.Println("Starting otel collector...")
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting otel collector:", err)
		return err
	}
	fmt.Println("Otel collector is started.")

	// Get the process ID
	pid := cmd.Process.Pid
	fmt.Println("Process ID:", pid)
	c.sync(true, &pid)

	return nil
}

func (c *Collector) Stop() error {
	// Get process ID
	pid := c.getPid()
	fmt.Println("PID: ", pid)

	// Find the process by its ID
	fmt.Println("Finding process...")
	process, err := os.FindProcess(*c.runnerSynchronizer.pid)
	if err != nil {
		fmt.Println("Process is not found: ", err)
		return err
	}
	fmt.Println("Process is found:.")

	// Send SIGTERM signal to the process
	fmt.Println("Stopping process...")
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Process could not be stopped: ", err)
		return err
	}
	fmt.Println("Process is stopped.")

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
