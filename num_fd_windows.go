// Fortio CLI util: number of open filedescriptor.
//
// (c) 2023 Fortio Authors
// See LICENSE

//go:build windows
// +build windows

package scli // import "fortio.org/scli"

import (
	"fortio.org/log"
	"golang.org/x/sys/windows"
	"unsafe"
)

func GetSystemInformation() (*windows.SYSTEM_PROCESS_INFORMATION, error) {
	// Sorta similar to https://go.googlesource.com/sys.git/+/master/windows/svc/security.go#83
	var systemProcessInfo *windows.SYSTEM_PROCESS_INFORMATION
	for infoSize := uint32((unsafe.Sizeof(*systemProcessInfo) + unsafe.Sizeof(uintptr(0))) * 1024); ; {
		systemProcessInfo = (*windows.SYSTEM_PROCESS_INFORMATION)(unsafe.Pointer(&make([]byte, infoSize)[0]))
		err := windows.NtQuerySystemInformation(windows.SystemProcessInformation, unsafe.Pointer(systemProcessInfo), infoSize, &infoSize)
		if err == nil {
			break
		} else if err != windows.STATUS_INFO_LENGTH_MISMATCH {
			return nil, err
		}
	}
	if log.LogDebug() {
		log.Debugf("GetSystemInformation: %#v", systemProcessInfo)
	}
	return systemProcessInfo, nil
}

func NumFD() int {
	systemProcessInfo, err := GetSystemInformation()
	if err != nil {
		log.Errf("GetSystemInformation failed: %v", err)
		return -1
	}
	return int(systemProcessInfo.HandleCount)
}
