package data

import "context"

type DataCenter interface {
	Range() (low uint64, high uint64)
	Get(ctx context.Context, ids ...uint64) ([]*File, error)
	Put(ctx context.Context, files ...*File) error
}

type File struct {
	Id   uint64 `json:"id"`
	Data []byte `json:"data"`
}
