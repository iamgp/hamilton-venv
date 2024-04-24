/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func addFilesToEnv(envName string, files []string) {
	envPath, _ := hvenvDir(envName)

	for _ , file := range files {
		file = strings.TrimSpace(file)

		absPath, err := filepath.Abs(file)
		destPath := filepath.Join(envPath, file)
		check(err)

		err = os.MkdirAll(filepath.Dir(destPath), 0755)
		check(err)

		_, err = copy(absPath, destPath)
		check(err)
	}
}

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "A brief description of your command",
	Long: ``,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
		addFilesToEnv(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
