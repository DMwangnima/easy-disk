package local

import (
	"context"
	"fmt"
	"github.com/DMwangnima/easy-disk/data/storage"
	"math/rand"
	"testing"
	"time"
)

const (
	TestByte = '1'
	TestFileSize = 4096
	TestMaxBlockNum = 1024
	TestMaxNum = 32768
    ChunkPath = "/Users/wangyuxuan/net-disk/data-server-test-dir/chunk"
    PoolPath = "/Users/wangyuxuan/net-disk/data-server-test-dir/pool"
    RosePath = "/Users/wangyuxuan/net-disk/data-server-test-dir/rose"
)

type Range struct {
	low  uint64
	high uint64
}

type Suite struct {
	// storage所能支持的最大id
	allSize   int
	// 每个block的大小，默认4kb
	blockSize int
	// 存储事先初始化好的字节切片，key为block数，从1到TestMaxBLockNum
	bufMap    map[int][]byte
}

func NewSuite(allSize, blockSize int) *Suite {
	bufMap := make(map[int][]byte)
	// 事先初始化
	for i := 1; i <= TestMaxBlockNum; i *= 2 {
		bufMap[i] = make([]byte, i*blockSize)
		for j := 0; j < i*blockSize; j++ {
			bufMap[i][j] = TestByte
		}
	}
	return &Suite{
		allSize:   allSize,
		blockSize: blockSize,
		bufMap:    bufMap,
	}
}

// ProduceGetReq: 传入序号以及需要的block数，返回该序号所指代的范围
// size: 每次取的block数
// num: 每次取的序号
func (suite *Suite) ProduceGetReq(num, size int) Range {
	rangeNum := suite.allSize / size
	low := num % rangeNum * size
	high := low + size - 1
	return Range{
		low:  uint64(low),
		high: uint64(high),
	}
}

// ProducePutReq: 传入序号以及需要的block数，返回该序号所指代的范围以及要传输的字节切片
// size: 每次取的block数
// num: 每次取的序号
func (suite *Suite) ProducePutReq(num, size int) (Range, []byte) {
	rang := suite.ProduceGetReq(num, size)
	buf := suite.bufMap[size]
	return rang, buf
}

var (
	//sop = NewStorageObjectPool(PoolPath, TestMaxNum, 10000, TestFileSize, 0.5)
	suite = NewSuite(TestMaxNum, TestFileSize)
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func basicStoragePut(b *testing.B, ctx context.Context, store storage.Storage, suite *Suite, seq, blockNum int) {
	req, buf := suite.ProducePutReq(seq, blockNum)
	if err := store.Put(ctx, &storage.Transfer{
		Low:  req.low,
		High: req.high,
		Data: buf,
	}); err != nil {
		b.Log(err)
	}
}

func basicStorageGet(b *testing.B, ctx context.Context, store storage.Storage, suite *Suite, seq, blockNum int) {
	rang := suite.ProduceGetReq(seq, blockNum)
	if _, err := store.Get(ctx, rang.low, rang.high); err != nil {
		b.Log(err)
	}
}

func BenchmarkStorageChunk(b *testing.B) {
	ctx := context.Background()
	sc, err := NewStorageChunk(ChunkPath, TestFileSize, 1024*1024, TestMaxNum)
    if err != nil {
    	b.Fatal(err)
	}
	// Put
	for blockNum := 1; blockNum <= TestMaxBlockNum; blockNum *= 2 {
		b.Run(fmt.Sprintf("sequential put%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
                basicStoragePut(b, ctx, sc, suite, i, blockNum)
			}
		})
		b.Run(fmt.Sprintf("random put%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
                basicStoragePut(b, ctx, sc, suite, rand.Int(), blockNum)
			}
		})
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				basicStoragePut(b, ctx, sc, suite, rand.Int(), blockNum)
			}
		})
	}
    // Get
	for blockNum := 1; blockNum <= TestMaxBlockNum; blockNum *= 2 {
		b.Run(fmt.Sprintf("sequential get%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStorageGet(b, ctx, sc, suite, i, blockNum)
			}
		})
		b.Run(fmt.Sprintf("random get%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStorageGet(b, ctx, sc, suite, rand.Int(), blockNum)
			}
		})
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				basicStorageGet(b, ctx, sc, suite, rand.Int(), blockNum)
			}
		})
	}
}

func BenchmarkStorageRose(b *testing.B) {
	ctx := context.Background()
	sr, err := NewStorageRose(RosePath, TestMaxNum, TestFileSize)
	if err != nil {
		b.Fatal(err)
	}
	// Put
	for blockNum := 1; blockNum <= TestMaxBlockNum; blockNum *= 2 {
		b.Run(fmt.Sprintf("sequential put%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStoragePut(b, ctx, sr, suite, i, blockNum)
			}
		})
		b.Run(fmt.Sprintf("random put%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStoragePut(b, ctx, sr, suite, rand.Int(), blockNum)
			}
		})
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				basicStoragePut(b, ctx, sr, suite, rand.Int(), blockNum)
			}
		})
	}
	// Get
	for blockNum := 1; blockNum <= TestMaxBlockNum; blockNum *= 2 {
		b.Run(fmt.Sprintf("sequential get%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStorageGet(b, ctx, sr, suite, i, blockNum)
			}
		})
		b.Run(fmt.Sprintf("random get%d", blockNum), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				basicStorageGet(b, ctx, sr, suite, rand.Int(), blockNum)
			}
		})
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				basicStorageGet(b, ctx, sr, suite, rand.Int(), blockNum)
			}
		})
	}
}