package videos

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
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

func RunCommand(commandStr string, args []string) error {
	cmd := exec.Command(commandStr, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(stdOut)
	log.Println(bytes)
	if err != nil {
		return err
	}

	log.Printf("Waiting for command to finish...")
	log.Printf("Process id is %v", cmd.Process.Pid)
	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			log.Printf("Exit error is %+v, error code: %v\n", exitError, exitError.ExitCode())
			return exitError
		}
	}
	return nil
}
