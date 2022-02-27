package util

import (
	"github.com/edsrzf/mmap-go"
	"os"
)

const (
	RDONLY = mmap.RDONLY
	RDWR = mmap.RDWR
)

type MMap mmap.MMap

func Map(file *os.File, prot, flag int) (MMap, error) {
	m, err := mmap.Map(file, prot, flag)
	return MMap(m), err
}

func (m MMap) Flush() error {
	return mmap.MMap(m).Flush()
}