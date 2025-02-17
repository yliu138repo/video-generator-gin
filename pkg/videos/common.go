package videos

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/viper"
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

// Run command and get the result asynchronously
func RunCommand(commandStr string, args []string, getResult func(cmd *exec.Cmd, cmdErr error)) (int, error) {
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
		getResult(cmd, err)
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

type IP struct {
	Query string
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

// This ensures that returned IP are public IP
func GetIP2() string {
	req, err := http.Get("http://ip-api.com/json/")
	if err != nil {
		return err.Error()
	}
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		return err.Error()
	}

	var ip IP
	json.Unmarshal(body, &ip)

	return ip.Query
}

func GenerateDownloadFilePath(filePath string, ip string, protocol string) string {
	return fmt.Sprintf("%s://%s%s/api/v1/videos/download?file_path=%s",
		protocol, GetIP2(), viper.Get("PORT").(string), url.QueryEscape(filePath))
}
