package process

import (
	"fmt"
	"os/exec"
)

func TranscribeWithPythonStreamingRealtime(audioFile, lang string, onData func(string)) error {
	cmd := exec.Command("python3", "./transcribe.py", "-i", audioFile, "-m", lang)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("error getting stdout pipe: %v", err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("error getting stderr pipe: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start command: %v", err)
	}

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stdout.Read(buf)
			if n > 0 {
				onData(string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := stderr.Read(buf)
			if n > 0 {
				fmt.Printf("stderr: %s\n", string(buf[:n]))
			}
			if err != nil {
				break
			}
		}
	}()

	return cmd.Wait()
}
