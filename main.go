package main

import (
	"fmt"
	"os"

	"github.com/JayJamieson/envrun/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "envrun: %v\n", err)
		os.Exit(1)
	}
}
