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
	//GetId() uint64
	//GetUsing() int
	//AddUsing(delta int)
	//GetFlag() Flag
	//SetFlag(status Flag)
	//GetGeneration() uint64
	//SetGeneration(generation uint64)
	//GetSequence() uint64
	//SetSequence(sequence uint64)
}

type Stream interface {
	Consume() (Object, bool)
	Produce(obj Object)
	Error() string
}

type Flag uint8

const (
	FREE Flag = iota
	READ
	WRITE
)

type Transfer struct {
	Id   uint64 `json:"id"`
	Data []byte `json:"data"`
}
