package local

import (
	"container/list"
	"errors"
	"github.com/DMwangnima/easy-disk/data/storage"
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
	id    uint64
	file  *os.File
	using int
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

func (obj *Object) GetUsing() int {
	return obj.using
}

func (obj *Object) AddUsing(delta int) {
	obj.using += delta
}

// todo 性能测试
// todo 考虑系统调用read write的并发安全性

type ObjectPool struct {
	lock       sync.Mutex
	basePath   string
	maxSize    int
	size       int
	items      map[uint64]*list.Element
	lruList    *list.List
	evicts     map[uint64]*list.Element
	exitChan   chan struct{}
	signalChan chan struct{}
}

func NewObjectPool(basePath string, maxSize int) *ObjectPool {
	return &ObjectPool{
		basePath:   basePath,
		maxSize:    maxSize,
		items:      make(map[uint64]*list.Element),
		lruList:    list.New(),
		evicts:     make(map[uint64]*list.Element),
		exitChan:   make(chan struct{}),
		signalChan: make(chan struct{}),
	}
}

func (op *ObjectPool) generatePath(id uint64) string {
	return path.Join(op.basePath, strconv.FormatUint(id, 10))
}

func (op *ObjectPool) allocateOne(id uint64) storage.Object {
	op.lock.Lock()
	if item, ok := op.items[id]; ok {
		item.Value.(*Object).AddUsing(1)
		op.lruList.MoveToFront(item)
		if _, ok := op.evicts[id]; ok {
			delete(op.evicts, id)
		}
		op.lock.Unlock()
		return item.Value.(*Object)
	}
	op.lock.Unlock()

	for {
		op.lock.Lock()
		// 有可用空间
		if op.size < op.maxSize {
			// todo 检查err
			obj, _ := NewObject(id, op.generatePath(id))
			obj.AddUsing(1)
			item := op.lruList.PushFront(obj)
			op.items[id] = item
			op.lock.Unlock()
			return obj
		}
		op.lock.Unlock()
		// 无可用空间，阻塞
		select {
		case <-op.signalChan:
		case <-op.exitChan:
			return nil
		}
	}
}

func (op *ObjectPool) Allocate(ids ...uint64) ([]storage.Object, error) {
	if len(ids) > op.maxSize {
		return nil, errors.New("beyond the max size")
	}
	res := make([]storage.Object, len(ids))
	for i, id := range ids {
        obj := op.allocateOne(id)
		res[i] = obj
	}
	return res, nil
}
