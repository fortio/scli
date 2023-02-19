package main

import (
	"fortio.org/cli"
	"fortio.org/scli"
)

func main() {
	cli.MinArgs = 2
	cli.MaxArgs = 4
	if !scli.ServerMain() {
		// in reality in both case we'd start some actual server
		return
	}
	select {}
}
