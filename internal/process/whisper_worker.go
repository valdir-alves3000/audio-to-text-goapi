package process

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
)

type WhisperRequest struct {
	File string `json:"file"`
	Lang string `json:"lang"`
}

type WhisperWorker struct {
	cmd    *exec.Cmd
	stdin  io.WriteCloser
	stdout *bufio.Reader
	mutex  sync.Mutex
}

var (
	workerInstance *WhisperWorker
	once           sync.Once
)

func GetWhisperWorker() (*WhisperWorker, error) {
	var err error

	once.Do(func() {
		workerInstance, err = startWorker()
	})

	return workerInstance, err
}

func startWorker() (*WhisperWorker, error) {

	cmd := exec.Command("python3", "./transcribe_worker.py")

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	reader := bufio.NewReader(stdoutPipe)

	// Wait for READY
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line != "READY\n" {
		return nil, fmt.Errorf("worker failed to start")
	}

	return &WhisperWorker{
		cmd:    cmd,
		stdin:  stdin,
		stdout: reader,
	}, nil
}

func (w *WhisperWorker) Transcribe(
	audioFile string,
	lang string,
	onSegment func(string),
) error {

	w.mutex.Lock()
	defer w.mutex.Unlock()

	req := WhisperRequest{
		File: audioFile,
		Lang: lang,
	}

	payload, err := json.Marshal(req)
	if err != nil {
		return err
	}

	_, err = fmt.Fprintln(w.stdin, string(payload))
	if err != nil {
		return err
	}

	for {
		line, err := w.stdout.ReadString('\n')
		if err != nil {
			return err
		}

		line = line[:len(line)-1]

		if line == "__END__" {
			break
		}

		if len(line) > 6 && line[:6] == "ERROR:" {
			return fmt.Errorf(line)
		}

		onSegment(line)
	}

	return nil
}