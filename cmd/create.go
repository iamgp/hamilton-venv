/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"
)

func createEnv(envName string) {

	homeDir, err := os.UserHomeDir()
	check(err)
	envPath := filepath.Join(homeDir, ".hvenv", envName)

	err = os.MkdirAll(envPath, os.ModePerm)
	check(err)

	log.Info("Environment created successfully.")
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

		createEnv(args[0])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
}
