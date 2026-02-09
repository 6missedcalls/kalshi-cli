package main

import (
	"os"

	"github.com/6missedcalls/kalshi-cli/internal/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
