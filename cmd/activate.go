/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "A brief description of your command",
	Long: "",
	Run: func(cmd *cobra.Command, args []string) {
        envName := args[0]
				targetPath := cmd.Flag("target").Value.String()

        homeDir, err := os.UserHomeDir()
        if err != nil {
            log.Fatal(err)
        }

				if targetPath == "" {
					targetPath, err = os.Getwd()
					if err != nil {
						log.Fatal(err)
					}
				}

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

        // Modify the shell prompt to show the active environment name
        fmt.Println("export PS1='(" + envName + ") $PS1'")
	},
}

func init() {
	rootCmd.AddCommand(activateCmd)
}
