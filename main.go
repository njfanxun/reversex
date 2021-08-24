package main

import (
	"github.com/njfanxun/reversex/cmd"
	_ "github.com/njfanxun/reversex/lib"
	"os"
)

func main() {
	if err := cmd.ExecuteCommand(); err != nil {
		os.Exit(1)
	}
}
