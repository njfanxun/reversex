package main

import (
	"os"
	"reversex/cmd"
	_ "reversex/lib"
)

func main() {
	if err := cmd.ExecuteCommand(); err != nil {
		os.Exit(1)
	}
}
