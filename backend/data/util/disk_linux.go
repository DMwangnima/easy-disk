//go:build linux
// +build linux

package util

func DiskFree(path string) (free uint64, err error) {
	stat := new(syscall.Statfs_t)
	if err = syscall.Statfs(path, stat); err != nil {
		return 0, err
	}
	reservedBlocks := stat.Bfree - stat.Bavail
	total := (stat.Blocks - reservedBlocks) * uint64(stat.Frsize)
	free = stat.Bavail * uint64(stat.Frsize)
	if total < free {
		return 0, fmt.Errorf("total space %d < free space %d, fs corruption at %s, using fsck", total, free, path)
	}
	return free, nil
}