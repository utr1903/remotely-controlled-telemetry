package otelcollector

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

type Collector struct {
}

func New() *Collector {
	return &Collector{}
}

func (c *Collector) Run() {
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}

	// Build the full path to the "app" executable in the current directory
	appPath := filepath.Join(currentDir, "/bin/otelcol-contrib")

	// Create a new Cmd struct for the "app" executable with the argument
	cmd := exec.Command(appPath, "--config=./bin/otel-config.yaml")

	// Start the process
	if err := cmd.Start(); err != nil {
		fmt.Println("Error starting process:", err)
		return
	}

	// Get the process ID
	pid := cmd.Process.Pid
	fmt.Println("Process ID:", pid)

	// Wait for the process to finish or be interrupted
	if err := cmd.Wait(); err != nil {
		fmt.Println("Process finished with error:", err)
	}

	// Stop the process by sending SIGTERM
	if err := stopProcess(pid); err != nil {
		fmt.Println("Error stopping process:", err)
	}
}

func stopProcess(pid int) error {
	// Find the process by its ID
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}

	// Send SIGTERM signal to the process
	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return err
	}

	return nil
}
