package local

import (
	"context"
	"errors"
	"github.com/DMwangnima/easy-disk/data/storage"
	"github.com/DMwangnima/easy-disk/data/util"
	"os"
	"path"
	"strconv"
	"sync"
)

type Chunk struct {
	id   uint64
	file *os.File
	mu   sync.RWMutex
	mmap util.MMap
}

func NewChunk(chunkId uint64, basePath string, chunkSize uint64) (*Chunk, error) {
	newPath := path.Join(basePath, strconv.FormatUint(chunkId, 10))
	file, err := os.OpenFile(newPath, os.O_CREATE|os.O_RDWR, 0666)
	// todo 包装错误
	if err != nil {
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	// 该文件可能是已经创建好的，所以需要检查原有长度是否和chunkSize相同
	if info.Size() != 0 && uint64(info.Size()) != chunkSize {
		return nil, errors.New("size incompatible")
	}
	// 每次都要将该文件设置成entrySize大小，因为mmap并不能帮助扩大文件
	if err = file.Truncate(int64(chunkSize)); err != nil {
		return nil, err
	}
	m, err := util.Map(file, util.RDWR, 0)
	if err != nil {
		return nil, err
	}
	res := Chunk{
		id:   chunkId,
		file: file,
		mmap: m,
	}
	return &res, nil
}

// 确保buf的length与size相同
func (chunk *Chunk) Read(buf []byte, offset, size uint64) {
	chunk.mu.RLock()
	copy(buf, chunk.mmap[offset:offset+size])
	chunk.mu.Unlock()
}

func (chunk *Chunk) Write(buf []byte, offset, size uint64) {
	chunk.mu.Lock()
	copy(chunk.mmap[offset:offset+size], buf)
	chunk.mmap.Flush()
	chunk.mu.Unlock()
}

type task func()

type chunkFunc func(buf []byte, low, high, chunkIdx uint64) task

type StorageChunk struct {
	// default: 4kb
	blockSize uint64
	// todo: 选择一个更合适的大小
	// default: 1gb
	chunkSize     uint64
	blockPerChunk uint64
	// 总chunk数
	chunkNum      uint64
	// 总block数
	blockNum      uint64
	// key为chunk id
	chunks        map[uint64]*Chunk
}

// todo: 对传入的参数做调整
func NewStorageChunk(basePath string, blockSize, chunkSize, blockNum uint64) storage.Storage {
    blockPerChunk := chunkSize / blockSize
    chunkNum := blockNum / blockPerChunk
    chunkMap := make(map[uint64]*Chunk)
    for i := uint64(0); i < chunkNum; i++ {
    	newChunk, err := NewChunk(i, basePath, chunkSize)
    	if err != nil {
    		return nil
		}
		chunkMap[i] = newChunk
	}
	return &StorageChunk{
		blockSize:     blockSize,
		chunkSize:     chunkSize,
		blockPerChunk: blockPerChunk,
		chunkNum:      chunkNum,
		blockNum:      blockNum,
		chunks:        chunkMap,
	}
}

func (sc *StorageChunk) Get(ctx context.Context, low, high uint64) (*storage.Transfer, error) {
	if low >= sc.blockNum {
		return nil, ErrStorageChunkBeyond
	}
	buf := make([]byte, (high-low+1)*sc.blockSize)
	tasks := sc.divideRange(buf, low, high, sc.chunkReadTask)
	var wg sync.WaitGroup
	for _, tk := range tasks {
		go func(t task) {
			t()
			wg.Done()
		}(tk)
	}
	wg.Wait()
	return &storage.Transfer{
		Low:  low,
		High: high,
		Data: buf,
	}, nil
}

// 获得chunk应处理的读取任务
func (sc *StorageChunk) chunkReadTask(buf []byte, low, high, chunkIdx uint64) task {
	return func() {
		blocks := high - low + 1
		offset := (low - (chunkIdx * sc.blockPerChunk)) * sc.blockSize
		sc.chunks[chunkIdx].Read(buf, offset, blocks*sc.blockSize)
	}
}

func (sc *StorageChunk) Put(ctx context.Context, transfer *storage.Transfer) error {
	if transfer.Low >= sc.blockNum {
		return ErrStorageChunkBeyond
	}
    tasks := sc.divideRange(transfer.Data, transfer.Low, transfer.High, sc.chunkWriteTask)
	var wg sync.WaitGroup
    wg.Add(len(tasks))
	for _, tk := range tasks {
		go func(t task) {
			t()
			wg.Done()
		}(tk)
	}
	wg.Wait()
	return nil
}

// 获得chunk应该处理的写入任务
func (sc *StorageChunk) chunkWriteTask(buf []byte, low, high, chunkIdx uint64) task {
	return func() {
		blocks := high - low + 1
		offset := (low - (chunkIdx * sc.blockPerChunk)) * sc.blockSize
		sc.chunks[chunkIdx].Write(buf, offset, blocks*sc.blockSize)
	}
}

func (sc *StorageChunk) divideRange(buf []byte, low, high uint64, fn chunkFunc) []task {
	lowChunkIdx, highChunkIdx := low/sc.blockPerChunk, high/sc.blockPerChunk
	res := make([]task, highChunkIdx-lowChunkIdx+1)
	if lowChunkIdx == highChunkIdx {
		res[0] = fn(buf, low, high, lowChunkIdx)
		return res
	}
	startOffset, endOffset := low, (lowChunkIdx+1)*sc.blockPerChunk-1
	for i := lowChunkIdx; i < highChunkIdx; i++ {
		res[i-lowChunkIdx] = fn(buf[(startOffset-low)*sc.blockSize:(endOffset+1-low)*sc.blockSize], startOffset, endOffset, i)
		startOffset = endOffset + 1
		endOffset += sc.blockPerChunk
	}
	endOffset = high
	res[len(res)-1] = fn(buf[(startOffset-low)*sc.blockSize:(endOffset+1-low)*sc.blockSize], startOffset, endOffset, highChunkIdx)
	return res
}
