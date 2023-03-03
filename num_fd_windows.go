// Fortio CLI util: number of open filedescriptor.
//
// (c) 2023 Fortio Authors
// See LICENSE

//go:build windows
// +build windows

package scli // import "fortio.org/scli"

import (
	"syscall"
	"unsafe"

	"fortio.org/log"
	"golang.org/x/sys/windows"
)

var (
	modkernel32           = windows.NewLazySystemDLL("kernel32.dll")
	getProcessHandleCount = modkernel32.NewProc("GetProcessHandleCount")
)

func GetCurrentProcessHandleCount() int {

	hdl, err := windows.GetCurrentProcess()
	if err != nil {
		log.Errf("GetCurrentProcess failed: %v", err)
		return -1
	}
	count := uint32(0)
	ret, _, err := syscall.Syscall(getProcessHandleCount.Addr(), 2, uintptr(hdl), uintptr(unsafe.Pointer(&count)), 0)
	log.Debugf("GetProcessHandleCount = %v, %v : %v", ret, err, count)
	if ret == 0 {
		log.Errf("GetProcessHandleCount failed: %v", err)
		return -1
	}
	return int(count)
}

func NumFD() int {
	return GetCurrentProcessHandleCount()
}
