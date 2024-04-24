// package main

// import (
// 	"flag"
// 	"fmt"
// 	"io"
// 	"io/ioutil"
// 	"log"
// 	"os"
// 	"path/filepath"
// 	"strings"

// 	"github.com/urfave/cli/v2"
// )

// // ---------------------------------------------------------------------------------------
// // 1. Create a new environment
// // ---------------------------------------------------------------------------------------
// func createEnv(envName string, files []string) {
// 	envPath := filepath.Join(".hvenv", envName)
// 	err := os.MkdirAll(envPath, os.ModePerm)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		srcFile, err := os.Open(file)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer srcFile.Close()

// 		destFile, err := os.Create(filepath.Join(envPath, file))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer destFile.Close()

// 		_, err = io.Copy(destFile, srcFile)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func createEnvCmd(args []string) {
// 	createCmd := flag.NewFlagSet("create", flag.ExitOnError)
// 	createEnvName := createCmd.String("env", "", "Environment name")
// 	createFiles := createCmd.String("files", "", "Comma-separated list of files")
// 	createCmd.Parse(args)
// 	createEnv(*createEnvName, strings.Split(*createFiles, ","))
// }

// // ---------------------------------------------------------------------------------------
// // 2. Activate an environment
// // ---------------------------------------------------------------------------------------
// func activateEnv(envName string) {
// 	envPath := filepath.Join(".hvenv", envName)
// 	files, err := ioutil.ReadDir(envPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	for _, file := range files {
// 		srcFile, err := os.Open(filepath.Join(envPath, file.Name()))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer srcFile.Close()

// 		destFile, err := os.Create(file.Name())
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer destFile.Close()

// 		_, err = io.Copy(destFile, srcFile)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func activateEnvCmd(args []string) {
// 	activateCmd := flag.NewFlagSet("activate", flag.ExitOnError)
// 	activateEnvName := activateCmd.String("env", "", "Environment name")
// 	activateCmd.Parse(args)
// 	activateEnv(*activateEnvName)
// }

// // ---------------------------------------------------------------------------------------
// // 3. Deactivate an environment
// // ---------------------------------------------------------------------------------------

// func deactivateEnv(files []string) {
// 	for _, file := range files {
// 		err := os.Remove(file)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func deactivateEnvCmd(args []string) {
// 	deactivateCmd := flag.NewFlagSet("deactivate", flag.ExitOnError)
// 	deactivateFiles := deactivateCmd.String("files", "", "Comma-separated list of files")
// 	deactivateCmd.Parse(args)
// 	deactivateEnv(strings.Split(*deactivateFiles, ","))
// }

// // ---------------------------------------------------------------------------------------
// // 4. List available environments
// // ---------------------------------------------------------------------------------------

// func listEnvs() {
// 	envsPath := ".hvenv"
// 	if !checkHvenvExists() {
// 			fmt.Println(".hvenv directory does not exist")
// 			os.Exit(1)
// 	}

// 	files, err := ioutil.ReadDir(envsPath)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	fmt.Println("Available environments:")
// 	for _, file := range files {
// 		if file.IsDir() {
// 			fmt.Println(file.Name())
// 		}
// 	}
// }

// func listEnvsCmd(args []string) {
// 	listEnvs()
// }

// // --------------------------------------------------------------------------------------
// // 5. Add files to an environment
// // --------------------------------------------------------------------------------------
// func addFilesToEnv(envName string, files []string) {
// 	envPath := filepath.Join(".hvenv", envName)
// 	for _, file := range files {
// 		srcFile, err := os.Open(file)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer srcFile.Close()

// 		destFile, err := os.Create(filepath.Join(envPath, file))
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer destFile.Close()

// 		_, err = io.Copy(destFile, srcFile)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 	}
// }

// func addFilesToEnvCmd(args []string) {
// 	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
// 	addEnvName := addCmd.String("env", "", "Environment name")
// 	addFiles := addCmd.String("files", "", "Comma-separated list of files")
// 	addCmd.Parse(args)
// 	addFilesToEnv(*addEnvName, strings.Split(*addFiles, ","))
// }

// // ---------------------------------------------------------------------------------------
// // 6. Main CLI function
// // ---------------------------------------------------------------------------------------
// type CommandFunc func(args []string)

// var commands = map[string]CommandFunc{
// 	"create":     createEnvCmd,
// 	"activate":   activateEnvCmd,
// 	"deactivate": deactivateEnvCmd,
// 	"list":       listEnvsCmd,
// 	"add":        addFilesToEnvCmd,
// }

// func main() {
// 	    (&cli.App{}).Run(os.Args)
// }

// func main2() {
// 	if len(os.Args) < 2 {
// 		fmt.Println("expected 'create', 'activate', 'deactivate', 'add' or 'list' subcommands")
// 		os.Exit(1)
// 	}

// 	cmdFunc, ok := commands[os.Args[1]]
// 	if !ok {
// 		fmt.Println("expected 'create', 'activate', 'deactivate', 'add' or 'list' subcommands")
// 		os.Exit(1)
// 	}

// 	cmdFunc(os.Args[2:])
// }
