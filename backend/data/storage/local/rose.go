package local

import (
	"context"
	"fmt"
	"github.com/DMwangnima/easy-disk/data/storage"
	"github.com/roseduan/rosedb"
	rose_storage "github.com/roseduan/rosedb/storage"
)

type StorageRose struct {
	db      *rosedb.RoseDB
	valSize uint64
	keyNum  uint64
}

func NewStorageRose(basePath string, keyNum, valSize uint64) storage.Storage {
	config := rosedb.DefaultConfig()
	config.RwMethod = rose_storage.MMap
	config.IdxMode = rosedb.KeyOnlyMemMode
	config.DirPath = basePath
	config.Sync = true
	db, err := rosedb.Open(config)
	if err != nil {
		panic(err)
	}
	return &StorageRose{
		db:      db,
		valSize: valSize,
		keyNum:  keyNum,
	}
}

func (sr *StorageRose) Get(ctx context.Context, low, high uint64) (*storage.Transfer, error) {
	if low >= sr.keyNum {
		return nil, ErrStorageChunkBeyond
	}
	vals, err := sr.db.RangeScan(low, high)
	if err != nil {
		return nil, err
	}
	buf := make([]byte, (high-low+1)*sr.valSize)
    for idx, val := range vals {
    	i := uint64(idx)
    	copy(buf[i*sr.valSize:(i+1)*sr.valSize], val.([]byte))
	}
	return &storage.Transfer{
		Low:  low,
		High: high,
		Data: buf,
	}, nil
}

func (sr *StorageRose) Put(ctx context.Context, trans *storage.Transfer) error {
	if trans.Low >= sr.keyNum {
		return ErrStorageChunkBeyond
	}
	for start := trans.Low; start <= trans.High; start++ {
		if err := sr.db.Set(start, trans.Data[(start-trans.Low)*sr.valSize:(start+1-trans.Low)*sr.valSize]); err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
