package util

import (
	"os"
	"path/filepath"
	"strings"
)

var supportedFormats = map[string]bool{
	".aac":  true,
	".flac": true,
	".m4a": true,
	".mp3": true,
	".mp4": true,
	".ogg":  true,
	".opus": true,
	".wav": true,
	".avi": true,
	".flv": true,
	".mkv": true,
	".mov": true,
	".mpeg": true,
	".webm": true,
}

func IsSupportedFormat(ext string) bool {
	ext = strings.ToLower(ext)
	return supportedFormats[ext]
}

func GetTempDir() string {
	return filepath.Join(os.TempDir(), "audio-transcription")
}
