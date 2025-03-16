# Video Compression Daemon

A Go daemon that automatically compresses MP4 files dropped into a watched folder using FFmpeg.

## Features

- Watches a folder for new `.mp4` files
- Automatically compresses videos using FFmpeg with H.264 codec
- Moves processed files to a `done/` folder to prevent reprocessing
- Handles partial writes by checking file stability
- Runs compression tasks asynchronously
- Supports automatic startup on macOS login

## Prerequisites

1. Go 1.16 or later
2. FFmpeg installed on your system

For macOS, install FFmpeg using Homebrew:
```bash
brew install ffmpeg
```

## Installation

1. Clone this repository
2. Install dependencies:
```bash
go mod tidy
```

## Usage

1. Build the binary:
```bash
go build -o video-watcher
```

2. Run the daemon:
```bash
./video-watcher
```

Or run it in the background:
```bash
nohup ./video-watcher > video-watcher.log 2>&1 &
```

## Directory Structure

- `videos/` - Drop your MP4 files here for compression
- `compressed/` - Compressed videos are saved here
- `done/` - Original files are moved here after compression

## Auto-start on macOS Login (Recommended Setup)

The daemon can be configured to start automatically when you log in to your Mac:

1. Build the binary first (if not already done):
```bash
go build -o video-watcher
```

2. Get your current directory path:
```bash
pwd
```

3. Create the LaunchAgent directory:
```bash
mkdir -p ~/Library/LaunchAgents
```

4. Create the LaunchAgent plist file:
```bash
cat > ~/Library/LaunchAgents/com.video.watcher.plist << EOL
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
    <dict>
        <key>Label</key>
        <string>com.video.watcher</string>
        <key>ProgramArguments</key>
        <array>
            <string>REPLACE_WITH_YOUR_PATH/video-watcher</string>
        </array>
        <key>WorkingDirectory</key>
        <string>REPLACE_WITH_YOUR_PATH</string>
        <key>RunAtLoad</key>
        <true/>
        <key>KeepAlive</key>
        <true/>
        <key>StandardOutPath</key>
        <string>/tmp/video-watcher.log</string>
        <key>StandardErrorPath</key>
        <string>/tmp/video-watcher.log</string>
    </dict>
</plist>
EOL
```

Replace `REPLACE_WITH_YOUR_PATH` with your actual path from step 2.

5. Load the service:
```bash
launchctl load ~/Library/LaunchAgents/com.video.watcher.plist
```

### Managing the Service

- **Start the service:**
```bash
launchctl load ~/Library/LaunchAgents/com.video.watcher.plist
```

- **Stop the service:**
```bash
launchctl unload ~/Library/LaunchAgents/com.video.watcher.plist
```

- **Check if running:**
```bash
launchctl list | grep com.video.watcher
```

- **View logs:**
```bash
tail -f /tmp/video-watcher.log
```

## Monitoring

Check if the daemon is running:
```bash
ps aux | grep video-watcher
```

View real-time logs:
```bash
tail -f /tmp/video-watcher.log
```

## Memory Usage

The daemon is lightweight, using approximately:
- 5-10MB when idle
- More memory only during active compression (depends on video size)

## Troubleshooting

1. If the service doesn't start:
   - Check the log file: `cat /tmp/video-watcher.log`
   - Verify the paths in the plist file
   - Ensure the binary has execute permissions: `chmod +x video-watcher`

2. If videos aren't being compressed:
   - Check if FFmpeg is installed: `which ffmpeg`
   - Verify the `videos` directory exists
   - Check the log file for errors

## Contributing

Pull requests are welcome. For major changes, please open an issue first to discuss what you would like to change.