package envsetup

import (
	"bufio"
	"hvenv/pkg/fileops"
	"os"
	"path/filepath"
	"regexp"
)

func CreateEnvironment(envName string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	envPath := filepath.Join(homeDir, ".hvenv", envName)
	err = os.MkdirAll(envPath, 0755)
	if err != nil {
		return "", err
	}

	return envPath, nil
}

func ProcessFile(filePath, envPath string, processedFiles map[string]bool, copiedFiles map[string]string) error {
	if processedFiles[filePath] {
		return nil
	}
	processedFiles[filePath] = true

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	includeRegex := regexp.MustCompile(`#include\s+"([^"]+)"`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := includeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			includedFile := matches[1]
			fullPath, found := fileops.ResolvePath(includedFile, filePath)
			if found {
				err = fileops.CopyFile(fullPath, envPath, copiedFiles)
				if err != nil {
					return err
				}
				err = ProcessFile(fullPath, envPath, processedFiles, copiedFiles)
				if err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}
