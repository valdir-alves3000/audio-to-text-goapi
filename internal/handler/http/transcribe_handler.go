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
	util "github.com/valdir-alves3000/audio-to-text-goapi/internal/utils"
)

var modelMapping = map[string]bool{
	"en": true,
	"pt": true,
}

type transcribeHandler struct {
	worker *process.WhisperWorker
}

func NewTranscribeHandler(worker *process.WhisperWorker) gin.HandlerFunc {
	h := &transcribeHandler{
		worker: worker,
	}
	return h.handle
}

func (h *transcribeHandler) handle(c *gin.Context) {

	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Access-Control-Allow-Origin", "*")

	if c.GetHeader("Accept") != "text/event-stream" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Request must accept text/event-stream",
		})
		return
	}

	ctx := c.Request.Context()

	lang := c.PostForm("lang")
	if _, valid := modelMapping[lang]; !valid {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported language",
		})
		return
	}

	file, header, err := c.Request.FormFile("audio")
	if err != nil {
		log.Printf("Failed to get audio file: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to get audio file",
		})
		return
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(header.Filename))
	if !util.IsSupportedFormat(ext) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Unsupported format",
		})
		return
	}

	tempFile, err := process.CreateTempFile(file, ext)
	if err != nil {
		log.Printf("Failed to create temp file: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create temporary file",
		})
		return
	}
	defer os.Remove(tempFile.Name())
	defer tempFile.Close()

	hasAudio, err := process.HasAudioStream(tempFile.Name())
	if err != nil || !hasAudio {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "No audio stream found",
		})
		return
	}

	wavFile, err := process.ConvertToWAVIfNeeded(tempFile.Name(), ext)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to convert audio",
		})
		return
	}
	if wavFile != tempFile.Name() {
		defer os.Remove(wavFile)
	}

	chunks, err := process.SplitAudioInChunks(wavFile, 10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to split audio",
		})
		return
	}

	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Streaming not supported",
		})
		return
	}

	for _, part := range chunks {

		select {
		case <-ctx.Done():
			log.Println("Client disconnected")
			return
		default:
		}

		err := h.worker.Transcribe(part, lang, func(segment string) {

			select {
			case <-ctx.Done():
				return
			default:
				fmt.Fprintf(c.Writer, "data: %s\n\n", segment)
				flusher.Flush()
			}
		})

		if err != nil {
			fmt.Fprintf(c.Writer, "event: error\ndata: %s\n\n", err.Error())
			flusher.Flush()
			return
		}
	}

	fmt.Fprintf(c.Writer, "event: end\ndata: stream completed\n\n")
	flusher.Flush()
}