package main

import (
	"fmt"
	"github.com/meschbach/minecraft-overseer/cmd"
	"os"
)

func main()  {
	root := cmd.NewOverseerCLI()
	if err := root.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
