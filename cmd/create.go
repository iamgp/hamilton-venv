package cmd

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// Function to create a new environment
func createEnv(envName string) string {
	homeDir, err := os.UserHomeDir()
	checkError(err)
	envPath := filepath.Join(homeDir, ".hvenv", envName)

	err = os.MkdirAll(envPath, os.ModePerm)
	checkError(err)

	log.Info("Environment created successfully.")
	return envPath
}

// Function to copy a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, sourceFile)
	return err
}

// Function to extract file paths from the given file
func extractFilePaths(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var paths []string
	scanner := bufio.NewScanner(file)
	includeRegex := regexp.MustCompile(`#include\s+"([^"]+)"`)

	for scanner.Scan() {
		line := scanner.Text()
		matches := includeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			paths = append(paths, matches[1])
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return paths, nil
}

// Function to resolve the full path of an included file and track its source
func resolvePath(baseDir, includePath string) (string, string) {
	// Check if the path is already absolute
	if filepath.IsAbs(includePath) {
		return includePath, "absolute"
	}

	// Check if the file exists in the same directory as the source file
	sameDirPath := filepath.Join(baseDir, includePath)
	if _, err := os.Stat(sameDirPath); err == nil {
		return sameDirPath, "sameDir"
	}

	// Check if the file exists in the "C:\\Program Files (x86)\\HAMILTON\\Library" directory
	libraryPath := filepath.Join("C:\\Program Files (x86)\\HAMILTON\\Library", includePath)
	if _, err := os.Stat(libraryPath); err == nil {
		return libraryPath, "library"
	}

	// If the file is not found in either location, return the original path
	return includePath, "unknown"
}

// Function to find files with the same base name but different extensions
func findFilesWithExtensions(baseDir, baseName string) ([]string, error) {
	var files []string
	err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasPrefix(filepath.Base(path), baseName) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Creates a new environment.",
	Long:  "This command allows you to create a new environment. You can specify the name of the environment during creation.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		targetPath := createEnv(args[0])

		sourceFile := os.Args[3]

		// Extract file paths from the source file
		paths, err := extractFilePaths(sourceFile)
		if err != nil {
			log.Info("Error extracting file paths:", err)
			return
		}

		// Get the directory of the source file
		sourceDir := filepath.Dir(sourceFile)

		// Map to track the source of each file
		fileSources := make(map[string]string)

		// Copy each file to the target directory
		for _, path := range paths {
			// Resolve the full path of the included file and track its source
			fullPath, source := resolvePath(sourceDir, path)
			fileSources[fullPath] = source

			// Find files with the same base name but different extensions
			baseName := strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
			sameDirFiles, err := findFilesWithExtensions(sourceDir, baseName)
			if err != nil {
				log.Info("Error finding files with extensions:", err)
				return
			}
			libraryFiles, err := findFilesWithExtensions("C:\\Program Files (x86)\\HAMILTON\\Library", baseName)
			if err != nil {
				log.Info("Error finding files with extensions:", err)
				return
			}

			// Copy files from the same directory
			for _, file := range sameDirFiles {
				targetFilePath := filepath.Join(targetPath, filepath.Base(file))
				err := os.MkdirAll(filepath.Dir(targetFilePath), 0755)
				if err != nil {
					log.Error("Error creating target directory:", "error", err)
					return
				}
				err = copyFile(file, targetFilePath)
				if err != nil {
					log.Info("Error copying file:", err)
					return
				}
				log.Info("Copied file", "source", file, "target", targetFilePath)
			}

			// Copy files from the library directory
			for _, file := range libraryFiles {
				targetFilePath := filepath.Join(targetPath, filepath.Base(file))
				err := os.MkdirAll(filepath.Dir(targetFilePath), 0755)
				if err != nil {
					log.Error("Error creating target directory:", "error", err)
					return
				}
				err = copyFile(file, targetFilePath)
				if err != nil {
					log.Info("Error copying file:", err)
					return
				}
				log.Info("Copied file", "source", file, "target", targetFilePath)
			}
		}

		// Save the file sources to a metadata file
		metadataFile := filepath.Join(targetPath, "file_sources.txt")
		file, err := os.Create(metadataFile)
		if err != nil {
			log.Info("Error creating metadata file:", err)
			return
		}
		defer file.Close()

		for path, source := range fileSources {
			_, err := file.WriteString(fmt.Sprintf("%s: %s\n", path, source))
			if err != nil {
				log.Info("Error writing to metadata file:", err)
				return
			}
		}

		log.Info("File sources saved to", metadataFile)
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
