package local

import (
	"errors"
	"github.com/DMwangnima/easy-disk/data/storage"
	"github.com/DMwangnima/easy-disk/data/util"
	"io"
	"os"
	"path"
	"strconv"
	"sync"
)

type ObjectManager struct {
	basePath string
}

func NewObjectManager(basePath string) *ObjectManager {
	return &ObjectManager{basePath: basePath}
}

// 可能会有too many open files的问题
func (om *ObjectManager) Allocate(ids ...uint64) ([]storage.Object, error) {
	var err error
	if len(ids) == 0 {
		return nil, errors.New("empty id slice")
	}
	objs := make([]storage.Object, len(ids))
	for i, id := range ids {
		newPath := om.generatePath(id)
		objs[i], err = NewObject(id, newPath)
		if err != nil {
			panic(err)
		}
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

// todo 性能测试
// todo 考虑系统调用read write的并发安全性

type ObjectPool struct {
	lock     sync.Mutex
	basePath string
	maxSize  int
	lru      util.Lru
}

func NewObjectPool(basePath string, maxSize int) *ObjectPool {
	return &ObjectPool{
		basePath: basePath,
		maxSize:  maxSize,
		lru:      util.NewLru(maxSize),
	}
}

func (op *ObjectPool) Allocate(ids ...uint64) ([]storage.Object, error) {
	if len(ids) > op.maxSize {
		return nil, errors.New("beyond the max size")
	}
	var res []storage.Object
    op.lock.Lock()
	defer op.lock.Unlock()
	for _, id := range ids {
		val, ok := op.lru.Get(id)
		if !ok {
			obj, err := NewObject(id, op.generatePath(id))
			if err != nil {
				// 打日志
				// 释放资源
				// 包装err
				return res, err
			}
			op.lru.Add(id, obj)
			res = append(res, obj)
		}
		res = append(res, val.(*Object))
	}
	return res, nil
}

func (op *ObjectPool) generatePath(id uint64) string {
return path.Join(op.basePath, strconv.FormatUint(id, 10))
}
