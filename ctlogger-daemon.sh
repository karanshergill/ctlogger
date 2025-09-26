#!/bin/bash

# CTLogger Daemon Control Script
# Simple script to run CTLogger as a background daemon

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BINARY="$SCRIPT_DIR/ctlogger"
PIDFILE="/tmp/ctlogger.pid"
LOGFILE="/tmp/ctlogger.log"

# Function to check if daemon is running
is_running() {
    if [[ -f "$PIDFILE" ]]; then
        local pid=$(cat "$PIDFILE")
        if kill -0 "$pid" 2>/dev/null; then
            echo "$pid"
            return 0
        else
            rm -f "$PIDFILE"
            return 1
        fi
    fi
    return 1
}

# Function to start daemon
start_daemon() {
    local pid
    if pid=$(is_running); then
        echo "CTLogger daemon is already running (PID: $pid)"
        return 1
    fi

    echo "Starting CTLogger daemon..."

    # Start the process in background, detached from terminal
    nohup "$BINARY" "$@" > "$LOGFILE" 2>&1 &
    local daemon_pid=$!

    # Save PID
    echo "$daemon_pid" > "$PIDFILE"

    # Wait a moment to check if it started successfully
    sleep 2
    if kill -0 "$daemon_pid" 2>/dev/null; then
        echo "CTLogger daemon started successfully (PID: $daemon_pid)"
        echo "Logs: $LOGFILE"
        return 0
    else
        echo "Failed to start CTLogger daemon"
        rm -f "$PIDFILE"
        return 1
    fi
}

# Function to stop daemon
stop_daemon() {
    local pid
    if pid=$(is_running); then
        echo "Stopping CTLogger daemon (PID: $pid)..."
        kill "$pid"

        # Wait for graceful shutdown
        local count=0
        while kill -0 "$pid" 2>/dev/null && [[ $count -lt 10 ]]; do
            sleep 1
            ((count++))
        done

        # Force kill if still running
        if kill -0 "$pid" 2>/dev/null; then
            echo "Force killing daemon..."
            kill -9 "$pid"
        fi

        rm -f "$PIDFILE"
        echo "CTLogger daemon stopped"
        return 0
    else
        echo "CTLogger daemon is not running"
        return 1
    fi
}

# Function to show status
show_status() {
    local pid
    if pid=$(is_running); then
        echo "CTLogger daemon is running (PID: $pid)"
        echo "Log file: $LOGFILE"
        return 0
    else
        echo "CTLogger daemon is not running"
        return 1
    fi
}

# Function to show logs
show_logs() {
    if [[ -f "$LOGFILE" ]]; then
        tail -f "$LOGFILE"
    else
        echo "No log file found at $LOGFILE"
        return 1
    fi
}

# Main script logic
case "${1:-}" in
    start)
        shift
        start_daemon "$@"
        ;;
    stop)
        stop_daemon
        ;;
    restart)
        shift
        stop_daemon
        sleep 2
        start_daemon "$@"
        ;;
    status)
        show_status
        ;;
    logs)
        show_logs
        ;;
    *)
        echo "Usage: $0 {start|stop|restart|status|logs} [ctlogger-options...]"
        echo ""
        echo "Commands:"
        echo "  start [options]  - Start CTLogger daemon with optional parameters"
        echo "  stop            - Stop CTLogger daemon"
        echo "  restart [opts]  - Restart CTLogger daemon with optional parameters"
        echo "  status          - Show daemon status"
        echo "  logs            - Show and follow daemon logs"
        echo ""
        echo "Examples:"
        echo "  $0 start -r domains.txt -v"
        echo "  $0 start -r domains.txt -o output/"
        echo "  $0 status"
        echo "  $0 logs"
        exit 1
        ;;
esac