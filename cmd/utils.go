package cmd

import (
	"os"
)

func checkHvenvExists() bool {
	_, err := os.Stat(".hvenv")
	return !os.IsNotExist(err)
}
