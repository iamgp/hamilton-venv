package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/log"
)

func HvenvRootDir() (string, bool) {
	homeDir, err := os.UserHomeDir()
	Check(err)
	envPath := filepath.Join(homeDir, ".hvenv")

	_, err = os.Stat(envPath)
	return envPath, !os.IsNotExist(err)
}

func HvenvDir(envName string) (string, bool) {
	homeDir, err := os.UserHomeDir()
	Check(err)
	envPath := filepath.Join(homeDir, ".hvenv", envName)

	_, err = os.Stat(envPath)
	return envPath, !os.IsNotExist(err)
}

func Check(err error, messages ...string) {
	if err != nil {
		for _, message := range messages {
			log.Error(message)
		}
		log.Fatal(err)
		os.Exit(1)
	}
}

func Copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)

	log.SetTimeFormat(time.Kitchen)
	log.Info("Copied file: " + dst)
	return nBytes, err
}

// GetCurrentTimestamp returns the current time formatted as RFC3339.
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
