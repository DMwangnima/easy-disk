package local
//
//import (
//	"context"
//	"github.com/DMwangnima/easy-disk/data/storage"
//	"testing"
//	"time"
//)
//
//const (
//	fileSize = 4096
//	basePath = "/Users/wangyuxuan/net-disk/data-server-test-dir"
//)
//
//var store = NewStorage(fileSize, basePath)
//
//func TestStorage_Put(t *testing.T) {
//	start := time.Now()
//	fileNum := 10000
//	transfers := make([]*storage.Transfer, fileNum)
//	for i := 0; i < fileNum; i++ {
//		transfers[i] = &storage.Transfer{
//			Id:   uint64(i+1),
//			Data: make([]byte, fileSize),
//		}
//	}
//	if err := store.Put(context.Background(), transfers...); err != nil {
//		t.Log(err)
//		return
//	}
//	t.Logf("Put time elapsed: %dms", time.Since(start).Milliseconds())
//}
//
//func TestStorage_PutConcurrent(t *testing.T) {
//	start := time.Now()
//	fileNum := 10000
//	transfers := make([]*storage.Transfer, fileNum)
//	for i := 0; i < fileNum; i++ {
//		transfers[i] = &storage.Transfer{
//			Id:   uint64(i+10001),
//			Data: make([]byte, fileSize),
//		}
//	}
//	if err := store.PutConcurrent(context.Background(), transfers...); err != nil {
//		t.Log(err)
//		return
//	}
//	t.Logf("Put time elapsed: %dms", time.Since(start).Milliseconds())
//}
//
////func TestStorage_Get(t *testing.T) {
////    start := time.Now()
////
////}
