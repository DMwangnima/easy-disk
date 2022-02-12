package util

import (
	"github.com/hashicorp/golang-lru/simplelru"
)

// todo 考虑加入构造函数和析构函数
// todo 去除Lru内部的锁
type Lru interface {
	Add(key interface{}, val interface{}) bool
	Get(key interface{}) (val interface{}, ok bool)
	Contains(key interface{}) (ok bool)
	Peek(key interface{}) (val interface{}, ok bool)
	Remove(key interface{}) (ok bool)
	RemoveOldest() (key interface{}, val interface{}, ok bool)
	GetOldest() (key interface{}, val interface{}, ok bool)
	Keys() []interface{}
	Len() int
	Purge()
	Resize(newSize int) int
}

type lruImplement struct {
	*simplelru.LRU
}

func NewLru(size int) Lru {
	cache, _ := simplelru.NewLRU(size, nil)
	return &lruImplement{cache}
}

