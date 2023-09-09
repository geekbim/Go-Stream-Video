package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strconv"
)

func createHls(inputFile, outputDir string, segmentDuration int) error {
	// Create the output directory if it does not exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	// Create the HLS playlist and segment the video using ffmpeg
	ffmpegCmd := exec.Command(
		"ffmpeg",
		"-i", inputFile,
		"-profile:v", "baseline", //baseline profile is compatible with most devices
		"-level", "3.0",
		"-start_number", "0", //start numbering segment from 0
		"-hls_time", strconv.Itoa(segmentDuration), //duration of each segment in seconds
		"-hls_list_size", "0", //keep all segments in the playlist
		"-f", "hls",
		fmt.Sprintf("%s/playlist.m3u8", outputDir),
	)

	output, err := ffmpegCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to create HLS: %v\nOutput: %s", err, string(output))
	}

	return nil
}

// addHeaders will act as middleware to give us CORS support
func addHeaders(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		h.ServeHTTP(w, r)
	}
}

func main() {
	inputFile := "./sample.mp4"
	outputDir := "./video"
	segmentDuration := 10 //duration of each segment in seconds

	if err := createHls(inputFile, outputDir, segmentDuration); err != nil {
		log.Fatalf("Error createing HLS: %v", err)
	}

	log.Println("HLS created successfully")

	port := 8080

	http.Handle("/", addHeaders(http.FileServer(http.Dir(outputDir))))
	fmt.Printf("Starting server on %v\n", port)
	log.Printf("Serving %s on HTTP port: %v\n", outputDir, port)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%v", port), nil))
}
