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

		if err := createAudioChunk(inputPath, outfile, start, chunkDuration); err != nil {
			return nil, err
		}

		chunks = append(chunks, outfile)
	}

	return chunks, nil
}

func createAudioChunk(inputPath, outputPath, start string, duration int) error {
	cmd := exec.Command("ffmpeg", "-i", inputPath, "-ss", start, "-t", fmt.Sprintf("%d", duration), "-c", "copy", "-y", outputPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error creating chunk: %v\n%s", err, stderr.String())
	}
	return nil
}

func HasAudioStream(filePath string) (bool, error) {
	cmd := exec.Command("ffprobe",
		"-v", "error",
		"-select_streams", "a", // seleciona apenas áudio
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
