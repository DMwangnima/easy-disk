package local

import (
	"errors"
)

const (
	packageName = "storage:local"
	objectPoolName = "ObjectPool"
	StorageChunkName = "StorageChunk"
)

var (
	ErrObjectPoolFlag = PackageErr(objectPoolName, "read write incompatible")
	ErrObjectPoolExit = PackageErr(objectPoolName, "ObjectPool exits")
	ErrStorageChunkBeyond = PackageErr(StorageChunkName, "beyond the range")
)

func PackageErr(object string, errStr string) error {
	return errors.New(packageName+":"+object+":"+errStr)
}
