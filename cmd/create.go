/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"io"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

func createEnv(envName string, files []string) {
	envPath := filepath.Join(".hvenv", envName)
	err := os.MkdirAll(envPath, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {
		srcFile, err := os.Open(file)
		if err != nil {
			log.Fatal(err)
		}
		defer srcFile.Close()

		destFile, err := os.Create(filepath.Join(envPath, file))
		if err != nil {
			log.Fatal(err)
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}

		createEnv(args[0], args[1:])
	},
}

func init() {
	rootCmd.AddCommand(createCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// createCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	createCmd.Flags().BoolP("files", "f", false, "Files to include")
}
