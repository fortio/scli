// Fortio CLI util: number of open filedescriptor.
//
// (c) 2023 Fortio Authors
// See LICENSE

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
	defer f.Close()
	names, err := f.Readdirnames(-1)
	if err != nil {
		log.Errf("Unable to read %s: %v", dir, err)
		return -1
	}
	if log.LogDebug() {
		log.Debugf("Found %d names in %s: %v", len(names), dir, names)
	}
	return len(names) - 3 // -3 for . and .. and the dir we just opened
}

func NumFD() int {
	switch runtime.GOOS {
	case "windows":
		log.Errf("NumFD not (yet) implemented on windows")
		return -1
	case "darwin":
		return countDir("/dev/fd")
	default:
		// assume everyone else has a /proc/self/fd
		return countDir("/proc/self/fd")
	}
}
