package main

import (
	"fmt"
	"os"

	"github.com/user/envchain/internal/cli"
)

func main() {
	if err := cli.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "envchain: %v\n", err)
		os.Exit(1)
	}
}
