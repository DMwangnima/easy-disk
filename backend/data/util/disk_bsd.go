//go:build darwin
// +build darwin

package util

import (
	"fmt"
	"syscall"
)

func DiskFree(path string) (free uint64, err error) {
	stat := new(syscall.Statfs_t)
	if err = syscall.Statfs(path, stat); err != nil {
		return 0, err
	}
	reservedBlocks := stat.Bfree - stat.Bavail
	total := (stat.Blocks - reservedBlocks) * uint64(stat.Bsize)
	free = stat.Bavail * uint64(stat.Bsize)
	if total < free {
		return 0, fmt.Errorf("total space %d < free space %d, fs corruption at %s, using fsck", total, free, path)
	}
	return free, nil
}
