package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// PathExists checks if a given path exists and returns whether it's a file or a directory.
func PathExists(path string) (exists bool, isDir bool) {
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false, false
	}
	return true, info.IsDir()
}

// Check if the current working directory ends with the target directory
func PwdEndsWith(target string) (bool, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return false, err
	}
	return strings.HasSuffix(pwd, target), nil
}

// CreateDir creates a directory. If recursive is true, it creates parent directories as needed.
func CreateDir(path string, recursive bool) error {
	var err error
	if recursive {
		err = os.MkdirAll(path, 0755) // Creates parent directories if needed
	} else {
		err = os.Mkdir(path, 0755) // Creates a single directory
	}

	return err
}

// -----------------------------------------------------------------------------
// docker-compose
// -----------------------------------------------------------------------------
// IsComposeRunning checks if specific services from a compose file are running
func IsComposeRunning(composePath string, serviceNames ...string) (bool, error) {
	args := []string{"-f", composePath, "ps", "-q"}
	args = append(args, serviceNames...)

	cmd := exec.CommandContext(context.Background(), "docker-compose", args...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, fmt.Errorf("error checking running state: %v\nOutput: %s", err, output)
	}

	return len(strings.TrimSpace(string(output))) > 0, nil
}

func RunServiceOnce(composePath, service string) error {
	running, err := IsComposeRunning(composePath, service)
	if err != nil {
		return err
	}

	if running {
		return fmt.Errorf("service %s is already running", service)
	}

	cmd := exec.CommandContext(
		context.Background(),
		"docker-compose",
		"-f", composePath,
		"up", "-d",
	)

	if output, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("failed to start: %v\nOutput: %s", err, output)
	}

	return nil
}

// StartDockerCompose starts Docker Compose services
func StartDockerCompose(composePath string) (string, error) {
	cmd := exec.CommandContext(context.Background(),
		"docker-compose",
		"-f", composePath,
		"up", "-d",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error starting services: %v\nOutput: %s", err, output)
	}

	return string(output), nil
}

// StopDockerCompose stops Docker Compose services
func StopDockerCompose(composePath string) (string, error) {
	cmd := exec.CommandContext(context.Background(),
		"docker-compose",
		"-f", composePath,
		"down",
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("error stopping services: %v\nOutput: %s", err, output)
	}

	return string(output), nil
}
