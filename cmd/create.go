package cmd

import (
	"bufio"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new environment.",
	Long:  "This command allows you to create a new environment. You can specify the name of the environment during creation.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		// targetPath := createEnv(args[0])

		// sourceFile := os.Args[3]
		log.Info(os.Args)
		envName := os.Args[2]
		sourceFile := os.Args[3]

		// Create environment directory
		envPath, err := createEnvironment(envName)
		if err != nil {
			fmt.Println("Error creating environment:", err)
			return
		}

		// Copy the source file first
		err = copySourceFile(sourceFile, envPath)
		if err != nil {
			fmt.Println("Error copying source file:", err)
			return
		}
		processedFiles := make(map[string]bool)
		copiedFiles := make(map[string]string)

		// Process source file and its dependencies
		err = processFile(sourceFile, envPath, processedFiles, copiedFiles)
		if err != nil {
			fmt.Println("Error processing files:", err)
			return
		}

		fmt.Println("Environment setup complete.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func createEnvironment(envName string) (string, error) {
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

func processFile(filePath, envPath string, processedFiles map[string]bool, copiedFiles map[string]string) error {
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
			fullPath, found := resolvePath(includedFile, filePath)
			if found {
				err = copyFile(fullPath, envPath, processedFiles, copiedFiles)
				if err != nil {
					return err
				}
				err = processFile(fullPath, envPath, processedFiles, copiedFiles)
				if err != nil {
					return err
				}
			}
		}
	}

	return scanner.Err()
}

func resolvePath(includePath, currentPath string) (string, bool) {
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

func copySourceFile(srcPath, envPath string) error {
	var baseDir, subDir string
	if strings.Contains(srcPath, "Library") {
		baseDir = "C:\\Program Files (x86)\\HAMILTON\\Library"
		subDir = "Library"
	} else if strings.Contains(srcPath, "Methods") {
		baseDir = "C:\\Program Files (x86)\\HAMILTON\\Methods"
		subDir = "Methods"
	} else {
		// Default to Methods if not specified
		baseDir = "C:\\Program Files (x86)\\HAMILTON\\Methods"
		subDir = "Methods"
	}

	// Calculate the relative path from the base directory
	_, err := filepath.Rel(baseDir, srcPath)
	if err != nil {
		return err
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

	return nil
}

func copyFile(srcPath, envPath string, processedFiles map[string]bool, copiedFiles map[string]string) error {
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
	hash, err := hashFile(srcPath)
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
	err = ioutil.WriteFile(filepath.Join(envPath, "files.hv"), jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
func hashFile(filePath string) (string, error) {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	hasher := sha256.New()
	hasher.Write(file)
	hash := hex.EncodeToString(hasher.Sum(nil))

	return hash, nil
}
