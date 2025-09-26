# CTLogger Daemon Usage Guide

CTLogger now supports running as a background daemon, so it continues monitoring Certificate Transparency logs even after you close your terminal.

## 🚀 Quick Start

### Method 1: Shell Script (Recommended)
```bash
# Start daemon with domain filtering
./ctlogger-daemon.sh start -r domains.txt -v

# Check status
./ctlogger-daemon.sh status

# View live logs
./ctlogger-daemon.sh logs

# Stop daemon
./ctlogger-daemon.sh stop

# Restart with new options
./ctlogger-daemon.sh restart -r new_domains.txt -o output/
```

### Method 2: Built-in Daemon Mode
```bash
# Start as daemon
./ctlogger -daemon -r domains.txt -v

# The parent process will exit, child runs in background
# Check if running: ps aux | grep ctlogger
```

## 📋 Available Commands

### Shell Script Commands
- `start [options]` - Start CTLogger daemon with optional parameters
- `stop` - Stop the running daemon
- `restart [options]` - Stop and restart with new parameters
- `status` - Show daemon status and PID
- `logs` - Show and follow daemon logs in real-time

### Built-in Daemon Flags
- `-daemon` - Run as daemon in background
- `-pidfile <path>` - Specify PID file location (default: `/tmp/ctlogger.pid`)

## 💡 Usage Examples

### Basic Monitoring (All Domains)
```bash
./ctlogger-daemon.sh start
```

### Filtered Monitoring with Verbose Logging
```bash
./ctlogger-daemon.sh start -r my_domains.txt -v
```

### JSON Output to Files
```bash
./ctlogger-daemon.sh start -r domains.txt -o output_dir/ -j
```

### Monitor with File Watching
```bash
./ctlogger-daemon.sh start -r domains.txt -f -v
```

## 📁 File Locations

- **PID File**: `/tmp/ctlogger.pid`
- **Log File**: `/tmp/ctlogger.log`
- **Database**: `./ctlogs.db` (in working directory)

## 🔍 Monitoring the Daemon

### Check if Running
```bash
./ctlogger-daemon.sh status
```

### View Real-time Logs
```bash
./ctlogger-daemon.sh logs
# or
tail -f /tmp/ctlogger.log
```

### Check Database Growth
```bash
ls -lh ctlogs.db
```

## 🛠 Troubleshooting

### Daemon Won't Start
1. Check if already running: `./ctlogger-daemon.sh status`
2. Check logs: `cat /tmp/ctlogger.log`
3. Verify domain file exists: `ls -la your_domains.txt`

### Permission Issues
- Ensure script is executable: `chmod +x ctlogger-daemon.sh`
- Check write permissions for `/tmp/`

### Database Locked
- Stop daemon: `./ctlogger-daemon.sh stop`
- Remove database: `rm -f ctlogs.db`
- Restart daemon

## 🔧 Advanced Configuration

### Custom PID File Location
```bash
./ctlogger -daemon -pidfile /var/run/ctlogger.pid -r domains.txt
```

### Multiple Instances
Run multiple instances with different PID files:
```bash
# Instance 1
./ctlogger -daemon -pidfile /tmp/ctlogger1.pid -r domains1.txt

# Instance 2
./ctlogger -daemon -pidfile /tmp/ctlogger2.pid -r domains2.txt
```

## 📊 Performance Tips

- Use domain filtering (`-r`) to reduce resource usage
- Monitor database size regularly
- Use output directories (`-o`) for organized results
- Enable verbose mode (`-v`) only for debugging

## 🔄 Auto-restart on Reboot

Add to your crontab for auto-start on system reboot:
```bash
crontab -e
# Add this line:
@reboot /path/to/ctlogger-daemon.sh start -r /path/to/domains.txt
```

## 🎯 Key Benefits

✅ **Persistent Monitoring**: Continues running after terminal closure
✅ **Easy Management**: Simple start/stop/status commands
✅ **Real-time Logs**: Live log monitoring and file output
✅ **Graceful Shutdown**: Proper signal handling and cleanup
✅ **Multiple Output Modes**: stdout, files, database, JSON
✅ **Resource Efficient**: Proper daemon implementation

The daemon mode ensures your Certificate Transparency monitoring runs continuously for maximum domain discovery coverage.