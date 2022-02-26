package space

import "io"

type Manager interface {
	Allocate(inode uint64, size uint64) (shardings []Sharding, err error)
	Resume(inode uint64) error
}

type Sharding interface {
	io.Writer
	Information() (offset uint64, size uint64)
}