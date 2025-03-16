package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

const (
	watchDir           = "./videos"
	outputDir          = "./compressed"
	doneDir            = "./done"
	compressionQuality = "28" // Lower means better quality
)

// Ensure directories exist
func ensureDir(path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.MkdirAll(path, os.ModePerm)
	}
}

// Compress MP4 file using FFmpeg
func compressMP4(inputPath string) {
	filename := filepath.Base(inputPath)
	outputPath := filepath.Join(outputDir, filename)

	// Skip if already compressed
	if _, err := os.Stat(outputPath); err == nil {
		log.Printf("Skipping %s (already compressed)", filename)
		return
	}

	// Run FFmpeg compression
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-vcodec", "libx264", "-crf", compressionQuality, outputPath)
	err := cmd.Run()
	if err != nil {
		log.Printf("Error compressing %s: %v\n", inputPath, err)
		return
	}

	log.Printf("Compressed: %s -> %s\n", inputPath, outputPath)

	// Move original file to done/
	donePath := filepath.Join(doneDir, filename)
	err = os.Rename(inputPath, donePath)
	if err != nil {
		log.Printf("Error moving %s to done/: %v\n", filename, err)
	} else {
		log.Printf("Moved %s to done/\n", filename)
	}
}

// Check if file size is stable (ensuring full write before processing)
func waitForCompleteWrite(path string) bool {
	var prevSize int64 = -1
	for i := 0; i < 5; i++ { // Check for 5 seconds
		info, err := os.Stat(path)
		if err != nil {
			return false
		}
		if info.Size() == prevSize {
			return true // File is stable
		}
		prevSize = info.Size()
		time.Sleep(1 * time.Second)
	}
	return false
}

func main() {
	ensureDir(outputDir)
	ensureDir(doneDir)

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Add(watchDir)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Watching folder:", watchDir)

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Create == fsnotify.Create && filepath.Ext(event.Name) == ".mp4" {
				log.Println("Detected new MP4:", event.Name)

				// Wait until the file is fully written
				if waitForCompleteWrite(event.Name) {
					go compressMP4(event.Name)
				} else {
					log.Println("Skipping (file not stable):", event.Name)
				}
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			log.Println("Watcher error:", err)
		}
	}
}
