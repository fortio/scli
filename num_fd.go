// Fortio CLI util: number of open filedescriptor.
//
// (c) 2023 Fortio Authors
// See LICENSE

//go:build !windows
// +build !windows

package scli // import "fortio.org/scli"

import (
	"os"
	"runtime"

	"fortio.org/log"
)

func countDir(dir string) int {
	f, err := os.Open(dir)
	if err != nil {
		log.Errf("Unable to open %s: %v", dir, err)
		return -1
	}
	/* to run lsof between stage to debug that on macos /dev/fd gets opened twice somehow
	if log.LogDebug() {
		log.Debugf("Sleeping after open before Readdirnames")
		time.Sleep(30 * time.Second)
		log.Debugf("Done sleeping, calling Readdirnames")
	}
	*/
	names, err := f.Readdirnames(-1)
	if err != nil {
		log.Errf("Unable to read %s: %v", dir, err)
		f.Close()
		return -1
	}
	if log.LogDebug() {
		log.Debugf("Found %d names in %s: %v", len(names), dir, names)
		// time.Sleep(60 * time.Second)
		// log.Debugf("Done sleeping, closing dir")
	}
	f.Close()
	return len(names) - 1 // for the dir we just opened
}

// NumFD returns the number of open file descriptors (or -1 on error).
// On windows it returns the number of handles.
func NumFD() int {
	switch runtime.GOOS {
	case "windows":
		log.Fatalf("Shouldn't be reached (!windows build tag)")
		return -1 // not reached (also, not possible)
	case "darwin":
		return countDir("/dev/fd") - 1 // macos seems to open 2 fds to Readdirnames /dev/fd
	default:
		// assume everyone else has a /proc/self/fd
		return countDir("/proc/self/fd")
	}
}
