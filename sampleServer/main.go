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
	log.Infof("FD count 1s after start : %d", scli.NumFD())
	time.Sleep(20 * time.Second)
	log.Infof("FD count 20s later      : %d", scli.NumFD())
	// is it stable:
	for range 5 {
		log.Infof("FD count stability check: %d", scli.NumFD())
	}
	if !srvStarted {
		// in reality in both case we'd start some actual server
		return
	}
	log.Infof("Running until interrupted (ctrl-c)...")
	scli.UntilInterrupted()
	log.Infof("Normal exit")
}
