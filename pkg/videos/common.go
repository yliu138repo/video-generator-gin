package videos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

type ErrorMessage struct {
	ErrorMsg string `json:"errorMessage"`
}

func CheckPathExist(path string) error {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("file path %s does not exist", path)
		}
		return err
	}
	return nil
}

// This will only kick off the command but no waiting the command to finish...s
func RunCommandContextNoWait(ctx context.Context, commandStr string, args []string) (int, error) {
	cmd := exec.CommandContext(ctx, commandStr, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}
	defer stdOut.Close()

	err = cmd.Start()
	log.Printf("Starting Process id is %v ...", cmd.Process.Pid)
	if err != nil {
		return -1, err
	}
	bytes, err := io.ReadAll(stdOut)
	log.Println(bytes)
	if err != nil {
		return -1, err
	}

	return cmd.Process.Pid, nil
}

func RunCommandContext(ctx context.Context, commandStr string, args []string) error {
	cmd := exec.CommandContext(ctx, commandStr, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	log.Printf("Starting Process id is %v ...", cmd.Process.Pid)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(stdOut)
	log.Println(bytes)
	if err != nil {
		return err
	}

	log.Printf("Waiting for command to finish...")
	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit error is %+v, error code: %v\n", exitError, exitError.ExitCode())
			return exitError
		}
	}
	return nil
}

// return true if the file passed is an image
func checkImage(fileName string) bool {
	exts := []string{".jpg", ".jpeg", ".png", ".git"}
	fileExt := filepath.Ext(fileName)

	for _, ext := range exts {
		if fileExt == ext {
			return true
		}
	}
	return false
}
