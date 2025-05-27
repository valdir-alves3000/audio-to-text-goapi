package process

import (
	"io"
	"os"
)

func CreateTempFile(file io.Reader, ext string) (*os.File, error) {
	tempFile, err := os.CreateTemp("", "audio-*"+ext)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(tempFile, file); err != nil {
		tempFile.Close()
		os.Remove(tempFile.Name())
		return nil, err
	}
	return tempFile, nil
}

