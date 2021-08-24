package main

import (
	"os"

	"github.com/njfanxun/reversex/cmd"
	_ "github.com/njfanxun/reversex/lib"
)

func main() {
	if err := cmd.ExecuteCommand(); err != nil {
		os.Exit(1)
	}
}
