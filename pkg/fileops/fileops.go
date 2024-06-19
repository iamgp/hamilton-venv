package fileops

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func ResolvePath(includePath, currentPath string) (string, bool) {
	// Check if the path is absolute and exists
	if filepath.IsAbs(includePath) {
		if _, err := os.Stat(includePath); err == nil {
			return includePath, true
		}
	}

	// Resolve relative to the current file's directory
	baseDir := filepath.Dir(currentPath)
	resolvedPath := filepath.Join(baseDir, includePath)
	if _, err := os.Stat(resolvedPath); err == nil {
		return resolvedPath, true
	}

	// Additionally check against known base directories
	knownBaseDirs := []string{
		"C:\\Program Files (x86)\\HAMILTON\\Library",
		"C:\\Program Files (x86)\\HAMILTON\\Methods",
	}

	for _, baseDir := range knownBaseDirs {
		resolvedPath = filepath.Join(baseDir, includePath)
		if _, err := os.Stat(resolvedPath); err == nil {
			return resolvedPath, true
		}
	}

	return "", false
}

func CopyFile(srcPath, envPath string, copiedFiles map[string]string) error {
	var baseDir, subDir string
	knownBaseDirs := map[string]string{
		"C:\\Program Files (x86)\\HAMILTON\\Library": "Library",
		"C:\\Program Files (x86)\\HAMILTON\\Methods": "Methods",
	}

	// Determine the correct base directory and subdirectory
	for dir, sub := range knownBaseDirs {
		if strings.HasPrefix(srcPath, dir) {
			baseDir = dir
			subDir = sub
			break
		}
	}

	if baseDir == "" {
		// If no known base directory matches, default to a general directory
		baseDir = "C:\\Program Files (x86)\\HAMILTON\\"
		subDir = ""
	}

	// Get the directory of the source file
	srcDir := filepath.Dir(srcPath)

	// Find files with the same base name but different extensions
	baseName := strings.TrimSuffix(filepath.Base(srcPath), filepath.Ext(srcPath))
	files, err := filepath.Glob(filepath.Join(srcDir, baseName+".*"))
	if err != nil {
		return err
	}

	// Copy each file to the target directory
	for _, file := range files {
		relativeFile, err := filepath.Rel(baseDir, file)
		if err != nil {
			return err
		}
		targetPath := filepath.Join(envPath, subDir, relativeFile)

		// Ensure the target directory exists
		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return err
		}

		// Open the source file
		sourceFile, err := os.Open(file)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		// Create the target file
		destinationFile, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		// Copy the file content
		_, err = io.Copy(destinationFile, sourceFile)
		if err != nil {
			return err
		}
	}

	// Calculate the hash of the source file
	hash, err := HashFile(srcPath)
	if err != nil {
		return err
	}

	// Add the file and its hash to the map of copied files
	copiedFiles[strings.ReplaceAll(srcPath, "\\\\", "\\")] = hash

	// Convert the map to JSON
	jsonData, err := json.MarshalIndent(copiedFiles, "", "	")
	if err != nil {
		return err
	}

	// Write the JSON data to the files.hv file
	err = os.WriteFile(filepath.Join(envPath, "files"), jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
func HashFile(filePath string) (string, error) {
	file, err := os.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(file)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}
