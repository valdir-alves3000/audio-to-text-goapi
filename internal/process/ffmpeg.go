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
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("file input not found: %s", inputPath)
	}

	ext := filepath.Ext(inputPath)
	outputPath := strings.TrimSuffix(inputPath, ext) + ".wav"

	cmd := exec.Command("ffmpeg",
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
		return "", fmt.Errorf("error converting to WAV: %v\n%s", err, stderr.String())
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return "", fmt.Errorf("failed to create output WAV file")
	}

	return outputPath, nil
}

func SplitAudioInChunks(inputPath string, chunkDuration int) ([]string, error) {
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("input file not found: %s", inputPath)
	}

	cmd := exec.Command("ffprobe", "-v", "error", "-show_entries",
		"format=duration", "-of", "default=noprint_wrappers=1:nokey=1", inputPath)

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("error getting audio duration: %v\n%s", err, out.String())
	}

	durationStr := strings.TrimSpace(out.String())
	duration, err := strconv.ParseFloat(durationStr, 64)
	if err != nil {
		return nil, fmt.Errorf("error converting duration: %v", err)
	}

	var chunks []string
	base := strings.TrimSuffix(inputPath, filepath.Ext(inputPath))

	for i := 0.0; i < duration; i += float64(chunkDuration) {
		start := fmt.Sprintf("%.2f", i)
		outfile := fmt.Sprintf("%s_part_%02d.wav", base, int(i)/chunkDuration)

		cmd := exec.Command("ffmpeg", "-i", inputPath, "-ss", start, "-t", fmt.Sprintf("%d", chunkDuration), "-c", "copy", "-y", outfile)

		var chunkErr bytes.Buffer
		cmd.Stderr = &chunkErr

		if err := cmd.Run(); err != nil {
			return nil, fmt.Errorf("erro ao criar chunk: %v\n%s", err, chunkErr.String())
		}

		chunks = append(chunks, outfile)
	}

	return chunks, nil
}
