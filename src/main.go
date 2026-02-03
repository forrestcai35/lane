package main

import (
	"os"

	"github.com/forrestcai35/lane/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
