package local

import (
	"errors"
)

const (
	packageName = "storage:local"
	objectPoolName = "ObjectPool"
)

var (
	ErrObjectPoolFlag = PackageErr(objectPoolName, "read write incompatible")
	ErrObjectPoolExit = PackageErr(objectPoolName, "ObjectPool exits")
)

func PackageErr(object string, errStr string) error {
	return errors.New(packageName+":"+object+":"+errStr)
}
