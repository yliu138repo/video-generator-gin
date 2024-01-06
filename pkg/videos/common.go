package videos

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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
func RunCommandContextNoWait(ctx context.Context, commandStr string, args []string) int {
	// cmd := exec.CommandContext(ctx, commandStr, args...)
	// cmd.Stdin = os.Stdin
	// cmd.Stderr = os.Stderr
	// // stdOut, err := cmd.StdoutPipe()
	// // if err != nil {
	// // 	return -1, err
	// // }
	// // defer stdOut.Close()

	// err := cmd.Start()
	// log.Printf("Starting Process id is %v ...", cmd.Process.Pid)
	// if err != nil {
	// 	log.Printf("Failed to start process %v:  %+v", cmd.Process.Pid, err)
	// }
	// // bytes, _ := io.ReadAll(stdOut)
	// // log.Println(bytes)
	// go func(cmd *exec.Cmd) {
	// 	log.Printf("Waiting for command to finish...")
	// 	err = cmd.Wait()
	// 	if err != nil {
	// 		if exitError, ok := err.(*exec.ExitError); ok {
	// 			log.Printf("Exit error is %+v, error code: %v\n", exitError, exitError.ExitCode())
	// 		}
	// 	}
	// }(cmd)

	// return cmd.Process.Pid, nil
	cmd := exec.CommandContext(ctx, commandStr, args...)
	go func() {
		data, err := cmd.CombinedOutput()
		log.Println(data, err, "===")
	}()
	return -1
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
			log.Println(cmd.Args, "===")
			return exitError
		}
	}
	return nil
}

func RunCommand(commandStr string, args []string) (int, error) {
	cmd := exec.Command(commandStr, args...)
	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	stdOut, err := cmd.StdoutPipe()
	if err != nil {
		return -1, err
	}

	err = cmd.Start()
	log.Printf("Starting Process id is %v ...", cmd.Process.Pid)
	if err != nil {
		return -1, err
	}

	go func() {
		bytes, err := io.ReadAll(stdOut)
		log.Println(bytes)
		if err != nil {
			log.Printf("Failed to read the stdout pipeline: %+v\n", err)
		}

		log.Printf("Waiting for command to finish...")
		err = cmd.Wait()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				log.Printf("Exit error is %+v, error code: %v\n", exitError, exitError.ExitCode())
			}
		}
	}()
	return cmd.Process.Pid, nil
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

// Get preferred outbound ip of this machine
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}
