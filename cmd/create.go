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
	reader := bufio.NewReader(file)
	includeRegex := regexp.MustCompile(`#include\s+"([^"]+)"`)

	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		matches := includeRegex.FindStringSubmatch(line)
		if len(matches) > 1 {
			paths = append(paths, matches[1])
		}
		if err == io.EOF {
			break
		}
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

	// Check if the file exists in the "C:\\Program Files (x86)\\HAMILTON\\Methods" directory
	methodsPath := filepath.Join("C:\\Program Files (x86)\\HAMILTON\\Methods", includePath)
	if _, err := os.Stat(methodsPath); err == nil {
		return methodsPath, "methods"
	}

	// If the file is not found in any location, return the original path
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

// Function to process a file and its dependencies recursively
func processFile(filePath, targetPath string, processedFiles map[string]bool, fileSources map[string]string) error {
	if processedFiles[filePath] {
		return nil
	}
	processedFiles[filePath] = true

	// Extract file paths from the source file
	paths, err := extractFilePaths(filePath)
	if err != nil {
		return fmt.Errorf("error extracting file paths: %w", err)
	}

	// Get the directory of the source file
	sourceDir := filepath.Dir(filePath)

	// Copy each file to the target directory
	for _, path := range paths {
		// Resolve the full path of the included file and track its source
		fullPath, source := resolvePath(sourceDir, path)
		fileSources[fullPath] = source

		// Find files with the same base name but different extensions
		baseName := strings.TrimSuffix(filepath.Base(fullPath), filepath.Ext(fullPath))
		sameDirFiles, err := findFilesWithExtensions(sourceDir, baseName)
		if err != nil {
			return fmt.Errorf("error finding files with extensions: %w", err)
		}
		libraryFiles, err := findFilesWithExtensions("C:\\Program Files (x86)\\HAMILTON\\Library", baseName)
		if err != nil {
			return fmt.Errorf("error finding files with extensions: %w", err)
		}
		methodsFiles, err := findFilesWithExtensions("C:\\Program Files (x86)\\HAMILTON\\Methods", baseName)
		if err != nil {
			return fmt.Errorf("error finding files with extensions: %w", err)
		}

		// Copy files from the same directory
		for _, file := range sameDirFiles {
			if processedFiles[file] {
				continue
			}
			processedFiles[file] = true
			fileSources[file] = "sameDir"

			relativePath, err := filepath.Rel(sourceDir, file)
			if err != nil {
				return fmt.Errorf("error getting relative path: %w", err)
			}
			targetFilePath := filepath.Join(targetPath, "Methods", relativePath)
			err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
			if err != nil {
				return fmt.Errorf("error creating target directory: %w", err)
			}
			err = copyFile(file, targetFilePath)
			if err != nil {
				return fmt.Errorf("error copying file: %w", err)
			}
			log.Info("Copied file", "source", file, "target", targetFilePath)
			// Process the copied file recursively
			err = processFile(file, targetPath, processedFiles, fileSources)
			if err != nil {
				return err
			}
		}

		// Copy files from the library directory
		for _, file := range libraryFiles {
			if processedFiles[file] {
				continue
			}
			processedFiles[file] = true
			fileSources[file] = "library"

			relativePath, err := filepath.Rel("C:\\Program Files (x86)\\HAMILTON\\Library", file)
			if err != nil {
				return fmt.Errorf("error getting relative path: %w", err)
			}
			targetFilePath := filepath.Join(targetPath, "Library", relativePath)
			err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
			if err != nil {
				return fmt.Errorf("error creating target directory: %w", err)
			}
			err = copyFile(file, targetFilePath)
			if err != nil {
				return fmt.Errorf("error copying file: %w", err)
			}
			log.Info("Copied file", "source", file, "target", targetFilePath)
			// Process the copied file recursively
			err = processFile(file, targetPath, processedFiles, fileSources)
			if err != nil {
				return err
			}
		}

		// Copy files from the methods directory
		for _, file := range methodsFiles {
			if processedFiles[file] {
				continue
			}
			processedFiles[file] = true
			fileSources[file] = "methods"

			relativePath, err := filepath.Rel("C:\\Program Files (x86)\\HAMILTON\\Methods", file)
			if err != nil {
				return fmt.Errorf("error getting relative path: %w", err)
			}
			targetFilePath := filepath.Join(targetPath, "Methods", relativePath)
			err = os.MkdirAll(filepath.Dir(targetFilePath), 0755)
			if err != nil {
				return fmt.Errorf("error creating target directory: %w", err)
			}
			err = copyFile(file, targetFilePath)
			if err != nil {
				return fmt.Errorf("error copying file: %w", err)
			}
			log.Info("Copied file", "source", file, "target", targetFilePath)
			// Process the copied file recursively
			err = processFile(file, targetPath, processedFiles, fileSources)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

		// Map to track processed files
		processedFiles := make(map[string]bool)
		// Map to track the source of each file
		fileSources := make(map[string]string)

		// Process the source file and its dependencies recursively
		err := processFile(sourceFile, targetPath, processedFiles, fileSources)
		if err != nil {
			log.Fatal("Error processing files:", err)
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
		log.Info("Environment created successfully with all dependencies.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
