package cmd

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/spf13/cobra"

	"hvenv/pkg/utils"
)

// activateCmd represents the activate command
var activateCmd = &cobra.Command{
	Use:   "activate",
	Short: "Activates a specific environment.",
	Long:  "This command allows you to activate a specific environment. Once activated, all subsequent operations will be performed within this environment.",
	Run:   activateEnvironment,
}

func init() {
	rootCmd.AddCommand(activateCmd)
}

func activateEnvironment(cmd *cobra.Command, args []string) {
	if len(args) < 1 {
		log.Fatal("Environment name is required")
	}
	envName := args[0]

	homeDir, err := os.UserHomeDir()
	checkError(err)

	targetPath, err := os.Getwd()
	checkError(err)

	envPath := filepath.Join(homeDir, ".hvenv", envName)
	switchEnvironmentFiles(envPath, targetPath)

	currentEnvFile := filepath.Join(homeDir, ".hvenv", ".current_hvenv")
	err = os.WriteFile(currentEnvFile, []byte(targetPath), 0644)
	checkError(err, "Failed to switch environment")

	log.Info("Environment switched successfully.", "env", envName, "path", currentEnvFile)
}

func switchEnvironmentFiles(envPath, targetPath string) {
	err := filepath.Walk(envPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			destPath := filepath.Join(targetPath, strings.TrimPrefix(path, envPath))
			_, err = utils.Copy(path, destPath)
			if err != nil {
				return err
			}
		}
		return nil
	})
	checkError(err)
}

func checkError(err error, msg ...string) {
	if err != nil {
		if len(msg) > 0 {
			log.Fatal(msg[0], err)
		} else {
			log.Fatal(err)
		}
	}
}
