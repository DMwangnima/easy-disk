//go:build windows
// +build windows

package util

import (
	"os"
	"unsafe"
)

var (
	kernel32 = windows.NewLazySystemDLL("kernel32.dll")
	// GetDiskFreeSpaceEx - https://msdn.microsoft.com/en-us/library/windows/desktop/aa364937(v=vs.85).aspx
	GetDiskFreeSpaceEx = kernel32.NewProc("GetDiskFreeSpaceExW")
)

func DiskFree(path string) (free uint64, err error) {
	if _, err = os.Stat(path); err != nil {
		return 0, err
	}
	availableBytes := int64(0)
	totalBytes := int64(0)
	freeBytes := int64(0)
	_, _, _ = GetDiskFreeSpaceEx.Call(uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(path))),
		    uintptr(unsafe.Pointer(&availableBytes),
			uintptr(unsafe.Pointer(&totalBytes)),
			uintptr(unsafe.Pointer(&freeBytes))))
	if uint64(totalBytes) < uint64(freeBytes) {
		return 0, fmt.Errorf("total space %d < free space %d, fs corruption at %s, using fsck", totalBytes, freeBytes, path)
	}
	return uint64(freeBytes), nil
}
