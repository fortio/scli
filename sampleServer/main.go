package main

import (
	"log"
	"time"

	"fortio.org/cli"
	"fortio.org/scli"
)

func main() {
	cli.MinArgs = 2
	cli.MaxArgs = 4
	srvStarted := scli.ServerMain()
	time.Sleep(1 * time.Second)
	log.Printf("FD count: %d", scli.NumFD())
	if !srvStarted {
		// in reality in both case we'd start some actual server
		return
	}
	select {}
}
