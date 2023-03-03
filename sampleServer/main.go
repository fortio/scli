package main

import (
	"time"

	"fortio.org/cli"
	"fortio.org/log"
	"fortio.org/scli"
)

func main() {
	cli.MinArgs = 2
	cli.MaxArgs = 4
	srvStarted := scli.ServerMain()
	time.Sleep(1 * time.Second)
	log.Infof("FD count: %d", scli.NumFD())
	time.Sleep(20 * time.Second)
	log.Infof("FD count: %d", scli.NumFD())
	// is it stable:
	for i := 0; i < 10; i++ {
		log.Infof("FD count: %d", scli.NumFD())
	}
	if !srvStarted {
		// in reality in both case we'd start some actual server
		return
	}
	select {}
}
