package main

import (
	"os"
	"xorm-reversex/cmd"
	_ "xorm-reversex/lib"
)

func main() {
	if err := cmd.ExecuteCommand(); err != nil {
		os.Exit(1)
	}
}
