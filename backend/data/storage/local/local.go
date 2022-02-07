package local

import (
	"context"
	"github.com/DMwangnima/easy-disk/data/storage"
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

func (s *Storage) PutConcurrent(ctx context.Context, transfers ...*storage.Transfer) error {
	if len(transfers) == 0 {
		return nil
	}
	ids := make([]uint64, len(transfers))
	objs, err := s.om.Allocate(ids...)
	if err != nil {
		return err
	}
	defer s.om.Resume(objs...)
	errChan := make(chan codes, 10)
	for i, trans := range transfers {
		go func() {
			num, err := objs[i].Write(trans.Data)
			// todo 统一错误
			if err != nil || num != len(trans.Data) {
				errChan <- err
			}
		}()
	}
	// todo 统一错误处理
	return nil
}
