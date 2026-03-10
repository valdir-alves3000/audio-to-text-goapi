package process

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

func ConvertToWAVIfNeeded(inputFile, ext string) (string, error) {
	if ext == ".wav" {
		return inputFile, nil
	}
	return convertToWAV(inputFile)
}

func convertToWAV(inputPath string) (string, error) {
	fileInfo, err := os.Stat(inputPath)
	if err != nil {
		return "", fmt.Errorf("input file not found: %w", err)
	}
	if fileInfo.Size() == 0 {
		return "", fmt.Errorf("input file is empty")
	}

	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + ".wav"

	if _, err := exec.LookPath("ffmpeg"); err != nil {
		return "", fmt.Errorf("ffmpeg not found in PATH: %w", err)
	}
	cmd := exec.Command("ffmpeg",
		"-v", "error", 
		"-i", inputPath, 
		"-ar", "16000", 
		"-ac", "1", 
		"-acodec", "pcm_s16le",
		"-y", 
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("error converting to WAV: %w\nDetails FFmpeg: %s", err, stderr.String())
	}

	if stat, err := os.Stat(outputPath); err != nil || stat.Size() == 0 {
		return "", fmt.Errorf("failed to create output WAV file (zero size or not created)")
	}

	return outputPath, nil
}

func SplitAudioInChunks(inputPath string, maxChunkDuration int) ([]string, error) {

	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file not found: %s", inputPath)
	}

	cmd := exec.Command(
		"ffmpeg",
		"-i", inputPath,
		"-af", "silencedetect=noise=-30dB:d=0.5",
		"-f", "null",
		"-",
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("silence detection failed: %v", err)
	}

	lines := strings.Split(stderr.String(), "\n")

	var silenceEnds []float64

	for _, line := range lines {
		if strings.Contains(line, "silence_end:") {
			parts := strings.Split(line, "silence_end:")
			if len(parts) < 2 {
				continue
			}

			value := strings.Fields(parts[1])[0]
			timestamp, err := strconv.ParseFloat(value, 64)
			if err == nil {
				silenceEnds = append(silenceEnds, timestamp)
			}
		}
	}

	if len(silenceEnds) == 0 {
		return []string{inputPath}, nil
	}

	var chunks []string
	base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))

	start := 0.0
	index := 0

	for _, end := range silenceEnds {

		if end-start < 1.0 {
			continue
		}

		if end-start > float64(maxChunkDuration) {
			end = start + float64(maxChunkDuration)
		}

		outfile := fmt.Sprintf("%s_part_%02d.wav", base, index)

		err := createAudioChunk(
			inputPath,
			outfile,
			fmt.Sprintf("%.3f", start),
			fmt.Sprintf("%.3f", end-start),
		)

		if err != nil {
			return nil, err
		}

		chunks = append(chunks, outfile)

		start = end
		index++
	}

	return chunks, nil
}

func createAudioChunk(inputPath, outputPath, start, duration string) error {

	cmd := exec.Command(
		"ffmpeg",
		"-v", "error",
		"-ss", start,
		"-i", inputPath,
		"-t", duration,
		"-ar", "16000",
		"-ac", "1",
		"-acodec", "pcm_s16le",
		"-y",
		outputPath,
	)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("chunk creation failed: %v\n%s", err, stderr.String())
	}

	return nil
}

func HasAudioStream(filePath string) (bool, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "a",
		"-show_entries", "stream=codec_type",
		"-of", "default=noprint_wrappers=1:nokey=1",
		filePath,
	)

	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("ffprobe failed: %w", err)
	}

	return len(output) > 0, nil
}
