/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activates a specific environment.",
	Long:  "This command allows you to activate a specific environment. Once activated, all subsequent operations will be performed within this environment.",
	Run: func(cmd *cobra.Command, args []string) {
		envName := args[0]
		// targetPath := cmd.Flag("target").Value.String()

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}

		targetPath, err := os.Getwd()
		check(err)

		envPath := filepath.Join(homeDir, ".hvenv", envName)

		// Switch out the files in the target directory with the files in the environment directory
		err = filepath.Walk(envPath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				destPath := filepath.Join(targetPath, strings.TrimPrefix(path, envPath))
				_, err = copy(path, destPath)
				if err != nil {
					return err
				}
			}

			return nil
		})
		if err != nil {
			log.Fatal(err)
		}

		currentEnvFile := filepath.Join(homeDir, ".hvenv", ".current_hvenv")

		err = os.WriteFile(currentEnvFile, []byte(targetPath), 0644)
		check(err, "Failed to switch environment")
		// err = os.WriteFile(currentEnvFile, []byte(envName), 0644)
		// check(err, "Failed to switch environment")
		log.Info("Environment switched successfully.", "env", envName, "path", currentEnvFile)
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
