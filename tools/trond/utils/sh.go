package utils

import (
	"os"
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
