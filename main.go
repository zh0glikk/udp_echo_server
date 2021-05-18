package main

import (
	"os"

	"udp_echo_server/internal/cli"
)

func main() {
	if !cli.Run(os.Args) {
		os.Exit(1)
	}
}
