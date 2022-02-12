package storage

import (
	"context"
	"io"
)

type Storage interface {
	Get(ctx context.Context, ids ...uint64) ([]*Transfer, error)
	Put(ctx context.Context, transfers ...*Transfer) error
	PutConcurrent(ctx context.Context, transfers ...*Transfer) error
}

type Object interface {
	io.ReadWriteCloser
	GetUsing() int
	AddUsing(delta int)
}

type Stream interface {
	Consume() (*Transfer, bool)
	Produce(trans *Transfer)
	Error() string
}

type Transfer struct {
	Id   uint64 `json:"id"`
	Data []byte `json:"data"`
}