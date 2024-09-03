package main

import (
	"binp/cli"
	"fmt"
	"os"
)

func main() {
	err := cli.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
