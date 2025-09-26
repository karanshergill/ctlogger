package main

import (
	"fmt"
	"os"

	"github.com/karanshergill/ctlogger/pkg/daemon"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ctlctl [start|stop|status|restart] [options...]")
		fmt.Println("  start   - Start CTLogger daemon")
		fmt.Println("  stop    - Stop CTLogger daemon")
		fmt.Println("  status  - Show daemon status")
		fmt.Println("  restart - Restart CTLogger daemon")
		fmt.Println("\nFor start command, pass ctlogger options after 'start'")
		fmt.Println("Example: ctlctl start -r domains.txt -v")
		os.Exit(1)
	}

	command := os.Args[1]
	pidFile := "/tmp/ctlogger.pid"

	switch command {
	case "start":
		// Check if already running
		running, pid, _ := daemon.IsRunning(pidFile)
		if running {
			fmt.Printf("CTLogger daemon is already running (PID: %d)\n", pid)
			os.Exit(1)
		}

		// Prepare args for ctlogger
		args := []string{"./ctlogger", "-daemon"}
		if len(os.Args) > 2 {
			args = append(args, os.Args[2:]...)
		}

		// Execute ctlogger with daemon flag
		err := daemon.Daemonize(pidFile)
		if err != nil {
			fmt.Printf("Failed to start daemon: %v\n", err)
			os.Exit(1)
		}

	case "stop":
		err := daemon.Stop(pidFile)
		if err != nil {
			fmt.Printf("Failed to stop daemon: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("CTLogger daemon stopped")

	case "status":
		status, err := daemon.Status(pidFile)
		if err != nil {
			fmt.Printf("Error checking status: %v\n", err)
			os.Exit(1)
		}
		fmt.Println(status)

	case "restart":
		// Stop first
		running, _, _ := daemon.IsRunning(pidFile)
		if running {
			err := daemon.Stop(pidFile)
			if err != nil {
				fmt.Printf("Failed to stop daemon: %v\n", err)
				os.Exit(1)
			}
			fmt.Println("CTLogger daemon stopped")
		}

		// Start again
		args := []string{"./ctlogger", "-daemon"}
		if len(os.Args) > 2 {
			args = append(args, os.Args[2:]...)
		}

		err := daemon.Daemonize(pidFile)
		if err != nil {
			fmt.Printf("Failed to start daemon: %v\n", err)
			os.Exit(1)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		fmt.Println("Available commands: start, stop, status, restart")
		os.Exit(1)
	}
}