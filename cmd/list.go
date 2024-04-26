/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func listEnvs() {
	rootDir, exists := hvenvRootDir()
	if !exists {
		log.Error(".hvenv directory does not exist. Please create an environment first.")
	}

	files, err := os.ReadDir(rootDir)
	if err != nil {
		log.Error(err)
	}

	log.Info("Available environments:")
	for _, file := range files {
		if file.IsDir() {
			log.Info(file.Name())
		}
	}
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists all available environments.",
	Long:  "This command allows you to list all available environments. The environments are displayed in the order they were created.",
	Run: func(cmd *cobra.Command, args []string) {
		listEnvs()
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
