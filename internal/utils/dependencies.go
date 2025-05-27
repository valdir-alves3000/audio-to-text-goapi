package util

import (
	"fmt"
	"os"
	"os/exec"
)

func CheckDependencies() error {
	required := []string{"ffmpeg", "python3"}
	for _, cmd := range required {
		if _, err := exec.LookPath(cmd); err != nil {
			return fmt.Errorf("%s not found in PATH", cmd)
		}
	}
	scriptPath := "./transcribe.py"
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("transcribe.py script not found in current directory")
	}
	if err := os.Chmod(scriptPath, 0755); err != nil {
		return fmt.Errorf("failed to make transcribe.py executable: %v", err)
	}
	return nil
}
