package daemon

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"syscall"
)

// Daemonize forks the current process and runs it in the background
func Daemonize(pidFile string) error {
	// Check if already running as daemon (second execution)
	if os.Getppid() == 1 {
		// We are the daemon process, write PID file
		return writePidFile(pidFile)
	}

	// First execution - fork and exit parent
	args := os.Args[:]
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = nil
	cmd.Stdout = nil
	cmd.Stderr = nil

	// Start the child process
	err := cmd.Start()
	if err != nil {
		return fmt.Errorf("failed to start daemon process: %v", err)
	}

	// Parent process exits
	fmt.Printf("CTLogger daemon started with PID %d\n", cmd.Process.Pid)
	os.Exit(0)
	return nil
}

// writePidFile writes the current process PID to the specified file
func writePidFile(pidFile string) error {
	pid := os.Getpid()
	pidStr := strconv.Itoa(pid)

	return os.WriteFile(pidFile, []byte(pidStr), 0644)
}

// IsRunning checks if a process with the PID from pidFile is running
func IsRunning(pidFile string) (bool, int, error) {
	data, err := os.ReadFile(pidFile)
	if err != nil {
		if os.IsNotExist(err) {
			return false, 0, nil
		}
		return false, 0, err
	}

	pid, err := strconv.Atoi(string(data))
	if err != nil {
		return false, 0, fmt.Errorf("invalid PID in file: %v", err)
	}

	// Check if process is running by sending signal 0
	process, err := os.FindProcess(pid)
	if err != nil {
		return false, pid, nil
	}

	err = process.Signal(syscall.Signal(0))
	if err != nil {
		return false, pid, nil
	}

	return true, pid, nil
}

// Stop stops the daemon process
func Stop(pidFile string) error {
	running, pid, err := IsRunning(pidFile)
	if err != nil {
		return err
	}

	if !running {
		return fmt.Errorf("daemon not running")
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		return fmt.Errorf("failed to find process: %v", err)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		return fmt.Errorf("failed to send SIGTERM: %v", err)
	}

	// Remove PID file
	os.Remove(pidFile)
	return nil
}

// Status returns the status of the daemon
func Status(pidFile string) (string, error) {
	running, pid, err := IsRunning(pidFile)
	if err != nil {
		return "", err
	}

	if running {
		return fmt.Sprintf("CTLogger daemon is running (PID: %d)", pid), nil
	}

	return "CTLogger daemon is not running", nil
}