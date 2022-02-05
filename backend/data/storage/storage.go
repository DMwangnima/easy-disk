package storage

import (
	"context"
	"io"
)

type Storage interface {
	Get(ctx context.Context, ids ...uint64) ([]*Transfer, error)
	Put(ctx context.Context, transfers ...*Transfer) error
}

type Object interface {
	io.ReadWriteCloser
}

type Transfer struct {
	Id   uint64 `json:"id"`
	Data []byte `json:"data"`
}