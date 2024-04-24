package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
)

func hvenvRootDir() (string, bool) {
	homeDir, err := os.UserHomeDir()
	check(err)
	envPath := filepath.Join(homeDir, ".hvenv")

	_, err = os.Stat(envPath)
	return envPath, !os.IsNotExist(err)
}

func hvenvDir(envName string) (string, bool) {
	homeDir, err := os.UserHomeDir()
	check(err)
	envPath := filepath.Join(homeDir, ".hvenv", envName)

	_, err = os.Stat(envPath)
	return envPath, !os.IsNotExist(err)
}

func check(err error, messages ...string) {
	if err != nil {
		for _, message := range messages {
			log.Error(message)
		}
		log.Fatal(err)
		os.Exit(1)
	}
}

func copy(src, dst string) (int64, error) {
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
	log.Info("File(s) copied.", "File",
	 strings.Join([]string{
			"source: " + src,
			"destination: " + dst,
			"size: " + fmt.Sprint(nBytes),
		}, "\n"),
	)
	return nBytes, err
}
