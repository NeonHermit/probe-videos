package analyzer

import (
	"encoding/json"
	"fmt"
	"github.com/NeonHermit/probe-videos/pkg/utils"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Stream struct {
	RFrameRate   string `json:"r_frame_rate"`
	NbReadFrames string `json:"nb_read_frames"`
	Duration     string `json:"duration"`
}

type FFProbeOutput struct {
	Streams []Stream `json:"streams"`
}

func RunAnalysis(dir string, csvFilePath string, logFn func(string)) {
	csvFile, err := os.Create(csvFilePath)
	if err != nil {
		logFn(fmt.Sprintf("Error creating CSV file: %v", err))
		return
	}
	defer csvFile.Close()

	header := "Filename, Frame Count, Duration (seconds), Frame Rate, Results\n"
	_, err = csvFile.WriteString(header)
	if err != nil {
		logFn(fmt.Sprintf("Error writing to CSV file: %v", err))
		return
	}

	files, err := os.ReadDir(dir)
	if err != nil {
		logFn(fmt.Sprintf("Error reading directory: %v", err))
		return
	}

	for _, file := range files {
		if file.IsDir() || !utils.IsSupportedFileType(file.Name()) {
			continue
		}

		fp := filepath.Join(dir, file.Name())
		if !utils.IsSupportedFileType(file.Name()) {
			logFn(fmt.Sprintf("Skipping unsupported file type: %s", file.Name()))
			continue
		}

		logFn(fmt.Sprintf("Processing: %s", file.Name()))

		cmd := exec.Command("ffprobe",
			"-show_entries", "stream=r_frame_rate,nb_read_frames,duration",
			"-select_streams", "v",
			"-count_frames",
			"-of", "json",
			"-threads", "3",
			"-v", "error", fp)

		output, err := cmd.Output()
		if err != nil {
			logFn(fmt.Sprintf("Error running ffprobe: %v", err))
			continue
		}

		var ffprobeOutput FFProbeOutput
		if err := json.Unmarshal(output, &ffprobeOutput); err != nil {
			logFn(fmt.Sprintf("Error parsing ffprobe output: %v", err))
			continue
		}

		if len(ffprobeOutput.Streams) > 0 {
			stream := ffprobeOutput.Streams[0]
			processStream(stream, fp, logFn, csvFile)
		}
	}
}

func processStream(stream Stream, fileName string, logFn func(string), csvFile *os.File) {
	frameRateFrac := strings.Split(stream.RFrameRate, "/")
	if len(frameRateFrac) == 2 {
		numerator, err1 := strconv.Atoi(frameRateFrac[0])
		denominator, err2 := strconv.Atoi(frameRateFrac[1])
		if err1 != nil || err2 != nil {
			logFn(fmt.Sprintf("Error converting frame rate to float for %s: %v %v", fileName, err1, err2))
			return
		}

		frameRateFloat := float64(numerator) / float64(denominator)
		frameCount, err := strconv.Atoi(stream.NbReadFrames)
		if err != nil {
			logFn(fmt.Sprintf("Error converting frame count to integer for %s: %v", fileName, err))
			return
		}

		results := float64(frameCount) / frameRateFloat
		duration, _ := strconv.ParseFloat(stream.Duration, 64)

		resultLine := fmt.Sprintf("%s, %d, %.2f, %.2f, %f\n",
			filepath.Base(fileName),
			frameCount, duration, frameRateFloat, results)

		_, err = csvFile.WriteString(resultLine)
		if err != nil {
			logFn(fmt.Sprintf("Error writing results for %s to CSV file: %v", fileName, err))
		} else {
			logFn(fmt.Sprintf("Results for %s written to CSV file.", fileName))
		}

		logFn(fmt.Sprintf("File: %s, Frame Count: %d, Duration: %.2f seconds, Frame Rate: %.2f, Results: %f",
			fileName, frameCount, duration, frameRateFloat, results))
	}
}
