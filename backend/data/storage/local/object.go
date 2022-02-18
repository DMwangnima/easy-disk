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
	"time"
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
		objs[i], err = NewObject(id, newPath, 0, 0)
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
	// generation:sequence to determine sequential order
	generation uint64
	sequence   uint64
	flag       storage.Flag
}

func NewObject(id uint64, path string, generation, sequence uint64) (*Object, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 777)
	if err != nil {
		return nil, err
	}
	obj := &Object{
		id:         id,
		file:       file,
		generation: generation,
		sequence:   sequence,
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

//func (obj *Object) GetId() uint64 {
//	return obj.id
//}
//
//func (obj *Object) GetUsing() int {
//	return obj.using
//}
//
//func (obj *Object) AddUsing(delta int) {
//	obj.using += delta
//}
//
//func (obj *Object) GetFlag() storage.Flag {
//	return obj.flag
//}
//
//func (obj *Object) SetFlag(flag storage.Flag) {
//	obj.flag = flag
//}
//
//func (obj *Object) GetGeneration() uint64 {
//	return obj.generation
//}
//
//func (obj *Object) SetGeneration(generation uint64) {
//	obj.generation = generation
//}
//
//func (obj *Object) GetSequence() uint64 {
//	return obj.sequence
//}
//
//func (obj *Object) SetSequence(sequence uint64) {
//	obj.sequence = sequence
//}

func (obj *Object) LessThan(than util.Item) bool {
	thanObj := than.(*Object)
	if obj.generation > thanObj.generation {
		return false
	}
	if obj.generation < thanObj.generation {
		return true
	}
	if obj.sequence > thanObj.sequence {
		return false
	}
	return true
}

// todo 性能测试
// todo 考虑系统调用read write的并发安全性

type ObjectPool struct {
	lock     sync.Mutex
	basePath string
	maxSize  int
	size     int
	// 关闭阈值
	evictThreshold int
	generation     uint64
	sequence       uint64
	items          map[uint64]*Object
	sortTree       util.RBTree
	exitChan       chan struct{}
	signalChan     chan struct{}
	resumeChan     chan struct{}
	resumeTimer    *time.Timer
}

// todo 考虑maxSize和evictRatio的默认值
func NewObjectPool(basePath string, maxSize int, evictRatio float64) *ObjectPool {
	return &ObjectPool{
		basePath:       basePath,
		maxSize:        maxSize,
		evictThreshold: int(float64(maxSize) * evictRatio),
		items:          make(map[uint64]*Object),
		sortTree:       util.NewRBTree(),
		exitChan:       make(chan struct{}),
		// todo 设置合理容量
		signalChan: make(chan struct{}, 10),
		resumeChan: make(chan struct{}, 10),
		// todo 设置回收时间
		resumeTimer: time.NewTimer(10 * time.Second),
	}
}

func (op *ObjectPool) generatePath(id uint64) string {
	return path.Join(op.basePath, strconv.FormatUint(id, 10))
}

func (op *ObjectPool) allocateOne(id uint64, flag storage.Flag) (storage.Object, error) {
	for {
		op.lock.Lock()
		newSequence := op.sequence + 1
		newGeneration := op.generation
		// 溢出
		if newSequence < op.sequence {
			newGeneration += 1
		}
		// 缓存命中
		if obj, ok := op.items[id]; ok {
			// 检查读写冲突
			if (flag == storage.WRITE && obj.flag != storage.FREE) || (flag == storage.READ && obj.flag == storage.WRITE) {
				op.lock.Unlock()
				return nil, ErrObjectPoolFlag
			}
			obj.using += 1
			obj.sequence = newSequence
			obj.generation = newGeneration
			obj.flag = flag

			op.sequence = newSequence
			op.generation = newGeneration
			op.sortTree.Delete(obj)
			op.lock.Unlock()
			return obj, nil
		}
		// 有可用空间，生成新Object
		// todo NewObject调用会造成一定阻塞，之后考虑状态机编程方式
		if op.size < op.maxSize {
			// todo 日志输出错误，考虑是否将该错误纳入ObjectPool中
			obj, err := NewObject(id, op.generatePath(id), newGeneration, newSequence)
			if err != nil {
				op.lock.Unlock()
				return nil, err
			}
			obj.using += 1
			obj.sequence = newSequence
			obj.generation = newGeneration
			obj.flag = flag

			op.items[id] = obj
			op.sequence = newSequence
			op.generation = newGeneration
			op.size += 1
			op.lock.Unlock()
			return obj, nil
		}
		op.lock.Unlock()
		// 无可用空间，阻塞，等待释放信号的到来
		select {
		case <-op.signalChan:
		case <-op.exitChan:
			return nil, ErrObjectPoolExit
		}
	}
}

func (op *ObjectPool) Allocate(flag storage.Flag, ids ...uint64) ([]storage.Object, error) {
	if len(ids) > op.maxSize {
		return nil, errors.New("beyond the max size")
	}
	res := make([]storage.Object, len(ids))
	for i, id := range ids {
		// todo 日志记录err
		// todo 检查err，并回退之前已经占用的Object
		obj, _ := op.allocateOne(id, flag)
		res[i] = obj
	}
	return res, nil
}

// AllocateStream 解决Allocate的队首阻塞问题，多个id有可能后面的id在缓存中，而第一个id阻塞了
// todo 性能测试
func (op *ObjectPool) AllocateStream(flag storage.Flag, ids ...uint64) storage.Stream {
	// todo 确定一个合理的chan size
    stream := NewStream(10)
    for _, id := range ids {
    	go func(oid uint64) {
    		obj, err := op.allocateOne(oid, flag)
    		// todo 日志记录
    		if err != nil {
				return
			}
			stream.Produce(obj)
		}(id)
	}
	return stream
}

func (op *ObjectPool) resumeOne(obj storage.Object) {
	op.lock.Lock()
	realObj := obj.(*Object)
	_, ok := op.items[realObj.id]
	if !ok {
		// todo 考虑panic或者记录日志
		op.lock.Unlock()
		return
	}
	realObj.using -= 1
	if realObj.using > 0 {
		op.lock.Unlock()
		return
	}
	op.sortTree.Insert(realObj)
	if op.sortTree.Len() >= op.evictThreshold {
		op.lock.Unlock()
		op.resumeChan <- struct{}{}
		return
	}
	op.lock.Unlock()
}

func (op *ObjectPool) Resume(objs ...storage.Object) {
	for _, obj := range objs {
		op.resumeOne(obj)
	}
}

func (op *ObjectPool) resumeDaemon() {
	go func() {
		for {
			select {
			case <-op.resumeChan:
				op.resumeTask()
			// 定时回收
			case <-op.resumeTimer.C:
				op.resumeTask()
			case <-op.exitChan:
				return
			}
		}
	}()
}

func (op *ObjectPool) resumeTask() {
	var evictSize int
	// 获得需要回收的size，暂定为排序树中待回收长度的一半，至少为1
	op.lock.Lock()
	treeLen := op.sortTree.Len()
	if treeLen <= 0 {
		op.lock.Unlock()
		return
	}
	evictSize = treeLen / 2
	if evictSize <= 0 {
		evictSize = 1
	}
	op.lock.Unlock()

	// 排序树中的待回收obj有可能在锁释放期间被重新使用，所以因检查此时树是否为空
	for i := 0; i < evictSize; i++ {
		// 每次释放一个obj后，即时释放锁，让排队的请求更快得到响应
		op.lock.Lock()
		obj := op.sortTree.DeleteMin().(*Object)
		if obj == nil {
			op.lock.Unlock()
			break
		}
		obj.Close()
		delete(op.items, obj.id)
		op.size -= 1
		op.lock.Unlock()

		select {
		case op.signalChan <- struct{}{}:
		case <-op.exitChan:
			return
		default:
		}
	}
}

// removeOne for testing to drawback some modifications.
func (op *ObjectPool) removeOne(id uint64) {
	os.Remove(op.generatePath(id))
}

func (op *ObjectPool) Exit() {
	close(op.exitChan)
}
