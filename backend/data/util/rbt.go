package util

import (
	"github.com/petar/GoLLRB/llrb"
)

// 项目内的代码实现该接口
type Item interface {
	LessThan(than Item) bool
}

// 适配器，用于调和llrb库的接口，使得项目内的代码只对util包内的接口产生依赖
type Adapter struct {
	Item Item
}

func (ada *Adapter) Less(than llrb.Item) bool {
	return ada.Item.LessThan(than.(*Adapter))
}

func (ada *Adapter) LessThan(than Item) bool {
	return ada.Item.LessThan(than)
}

type RBTree interface {
    DeleteMin() Item
    Delete(key Item) Item
    Insert(key Item)
    Len() int
}

func NewRBTree() RBTree {
	return &rbt{internal: llrb.New()}
}

type rbt struct {
	internal *llrb.LLRB
}

func (tree *rbt) DeleteMin() Item {
	item := tree.internal.DeleteMin()
	if item == nil {
		return nil
	}
	return item.(*Adapter).Item
}

func (tree *rbt) Delete(key Item) Item {
	item := tree.internal.Delete(&Adapter{Item: key})
	if item == nil {
		return nil
	}
	return item.(*Adapter).Item
}

func (tree *rbt) Insert(key Item) {
	tree.internal.ReplaceOrInsert(&Adapter{Item: key})
}

func (tree *rbt) Len() int {
	return tree.internal.Len()
}

