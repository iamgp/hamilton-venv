package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"hvenv/pkg/envsetup"
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

		envName := os.Args[2]
		sourceFile := os.Args[3]

		// Create environment directory
		envPath, err := envsetup.CreateEnvironment(envName)
		if err != nil {
			log.Info("Error creating environment:", err)
			return
		}

		// Copy the source file first
		err = copySourceFile(sourceFile, envPath)
		if err != nil {
			log.Info("Error copying source file:", err)
			return
		}

		// Process source file and its dependencies
		processedFiles := make(map[string]bool)
		copiedFiles := make(map[string]string)
		err = envsetup.ProcessFile(sourceFile, envPath, processedFiles, copiedFiles)
		if err != nil {
			log.Error("Error processing files:", "error", err)
			return
		}

		// Update info file with timestamp and updated tag
		infoFile := filepath.Join(envPath, "info")
		timestamp := time.Now().Format(time.RFC3339)
		content := fmt.Sprintf("updated: %s", timestamp)
		err = os.WriteFile(infoFile, []byte(content), 0644)
		if err != nil {
			log.Error("Error updating info file:", "error", err)
			return
		}

		log.Info("Environment setup complete.")
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
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
