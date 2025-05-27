package util

import (
	"os"
	"path/filepath"
	"strings"
)

var supportedFormats = map[string]bool{
	".wav": true,
	".mp3": true,
	".m4a": true,
	".mp4": true,
}

func IsSupportedFormat(ext string) bool {
	ext = strings.ToLower(ext)
	return supportedFormats[ext]
}

func GetTempDir() string {
	return filepath.Join(os.TempDir(), "audio-transcription")
}
