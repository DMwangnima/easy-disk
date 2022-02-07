package local

import (
	"errors"
	"github.com/DMwangnima/easy-disk/data/storage"
	"io"
	"os"
	"path"
	"strconv"
)

type ObjectManager struct {
    basePath string
}

func NewObjectManager(basePath string) *ObjectManager {
	return &ObjectManager{basePath: basePath}
}

func (om *ObjectManager) Allocate(ids ...uint64) ([]storage.Object, error) {
	var err error
	if len(ids) == 0 {
		return nil, errors.New("empty id slice")
	}
	objs := make([]storage.Object, len(ids))
	for i, id := range ids {
		newPath := om.generatePath(id)
		objs[i], err = NewObject(id, newPath)
	}
    return objs, err
}

func (om *ObjectManager) generatePath(id uint64) string {
	return path.Join(om.basePath, strconv.FormatUint(id, 10))
}

func (om *ObjectManager) Resume(objs ...storage.Object) {
	for _, obj := range objs {
		obj.Close()
	}
	return
}

type Object struct {
	id   uint64
	file *os.File
}

func NewObject(id uint64, path string) (storage.Object, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 777)
	if err != nil {
		return nil, err
	}
    obj := &Object{
		id:   id,
		file: file,
	}
	return obj, nil
}

// todo 考虑统一错误
func (obj *Object) Read(buf []byte) (n int, err error) {
	return io.ReadFull(obj.file, buf)
}

func (obj *Object) Write(buf []byte) (n int, err error) {
	return obj.file.Write(buf)
}

func (obj *Object) Close() error {
	return obj.file.Close()
}