package local

import (
	"context"
	"fmt"
	"github.com/DMwangnima/easy-disk/data/storage"
	"sync"
)

type Storage struct {
	om   *ObjectManager
	size int
}

func NewStorage(size int, basePath string) storage.Storage {
	om := NewObjectManager(basePath)
	return &Storage{
		om:   om,
		size: size,
	}
}

// sequential read
func (s *Storage) Get(ctx context.Context, ids ...uint64) ([]*storage.Transfer, error) {
	objs, err := s.om.Allocate(ids...)
	if err != nil {
		return nil, err
	}
	defer s.om.Resume(objs...)
	res := make([]*storage.Transfer, len(ids))
	for i, obj := range objs {
		// todo 后续考虑内存复用
		buf := make([]byte, s.size)
		num, err := obj.Read(buf)
		if err != nil || num != s.size {
			// 重复处理error
			return nil, err
		}
		res[i] = &storage.Transfer{
			Id:   ids[i],
			Data: buf,
		}
	}
	return res, nil
}

func (s *Storage) GetConcurrent(ctx context.Context, ids ...uint64) ([]*storage.Transfer, error) {
	objs, err := s.om.Allocate(ids...)
	if err != nil {
		return nil, err
	}
	defer s.om.Resume(objs...)

}

// sequential write
// 浪费了一些时间，需要进行优化
// todo 优化接口设计
func (s *Storage) Put(ctx context.Context, transfers ...*storage.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	ids := make([]uint64, len(transfers))
	for i, transfer := range transfers {
		ids[i] = transfer.Id
	}
	objs, err := s.om.Allocate(ids...)
	if err != nil {
		return err
	}
	defer s.om.Resume(objs...)
	// 顺序写
	for i, trans := range transfers {
		num, err := objs[i].Write(trans.Data)
		// todo 统一错误
		if err != nil || num != len(trans.Data) {
			return err
		}
	}
	return nil
}

// concurrent write
// 经测试比sequential write更快
// todo 调整goroutine量，因为磁盘io有限度，所以达到一定量就不再增长，反而会因为goroutine过多而增加性能损耗
func (s *Storage) PutConcurrent(ctx context.Context, transfers ...*storage.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	ids := make([]uint64, len(transfers))
	for i, transfer := range transfers {
		ids[i] = transfer.Id
	}
	objs, err := s.om.Allocate(ids...)
	if err != nil {
		return err
	}
	defer s.om.Resume(objs...)
	//errChan := make(chan error, 10)
	wg := sync.WaitGroup{}
	wg.Add(len(transfers))
	// 此处是使用循环变量的陷阱，必须要将trans复制一份
	for i, trans := range transfers {
		go func(trans *storage.Transfer) {
			defer wg.Done()
			num, err := objs[i].Write(trans.Data)
			// todo 统一错误
			if err != nil || num != len(trans.Data) {
				fmt.Println(err)
			}
		}(trans)
	}
	wg.Wait()
	// todo 统一错误处理
	return nil
}

func (s *Storage) Run() {
	go func() {

	}()
}
