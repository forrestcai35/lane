package main

import (
	"os"

	"github.com/forrest/lane/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
