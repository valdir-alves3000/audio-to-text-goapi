package http

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/process"
	"github.com/valdir-alves3000/audio-to-text-goapi/internal/utils"
)

var modelMapping = map[string]bool{
	"en-us": true,
	"pt":    true,
}

func TranscribeHandler(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	if c.GetHeader("Accept") != "text/event-stream" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Request must accept text/event-stream"})
		return
	}

	lang := c.PostForm("lang")
	if _, valid := modelMapping[lang]; !valid {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported language",
		})
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		log.Println("failed to get audio file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to get audio file"})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !util.IsSupportedFormat(ext) {
		log.Println("Unsupported format: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported format. Use WAV, MP3, M4A, or MP4",
		})
		return
	}

	tempFile, err := process.CreateTempFile(file, ext)
	if err != nil {
		log.Println("failed to create temporary file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create temporary file"})
		return
	}

	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	wavFile, err := process.ConvertToWAVIfNeeded(tempFile.Name(), ext)
	if err != nil {
		log.Println("failed to convert to WAV: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to convert to WAV"})
		return
	}
	if wavFile != tempFile.Name() {
		defer os.Remove(wavFile)
	}

	chunks, err := process.SplitAudioInChunks(wavFile, 2)
	if err != nil {
		log.Println("failed to split audio in chunks: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	done := make(chan bool)

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		log.Println("Streaming not supported")
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": "Streaming not supported",
		})
		return
	}

	go func() {
		defer close(done)

		ctx := c.Request.Context()
		for _, part := range chunks {
			select {
			case <-ctx.Done():
				return
			default:
				err := process.TranscribeWithPythonStreamingRealtime(part, lang, func(data string) {
					select {
					case <-ctx.Done():
						return // Skip write if context canceled
					default:
						fmt.Fprintf(c.Writer, "data: %s\n\n", strings.TrimSpace(data))
						flusher.Flush()
					}
				})

				if err != nil {
					select {
					case <-ctx.Done():
						return
					default:
						fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
						flusher.Flush()
					}
					return
				}
			}
		}

		select {
		case <-ctx.Done():
			return
		default:
			fmt.Fprintf(c.Writer, "event: end\ndata: stream completed\n\n")
			flusher.Flush()
		}
	}()

	select {
	case <-done:
	case <-c.Request.Context().Done():
	}
}
