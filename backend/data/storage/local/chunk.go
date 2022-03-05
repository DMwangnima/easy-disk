package local

import (
	"context"
	"errors"
	"fmt"
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
	chunk.mu.RUnlock()
}

func (chunk *Chunk) Write(buf []byte, offset, size uint64) {
	chunk.mu.Lock()
	copy(chunk.mmap[offset:offset+size], buf)
	chunk.mmap.Flush()
	chunk.mu.Unlock()
}

// bitMap: StorageChunk应该有知晓哪些序号已被使用的能力，考虑到有可能宕机，内存中的set或bitMap并不靠谱，因此使用磁盘文件较为妥当。
// 1.在当前场景中，假设磁盘容量最大为1TB，即2^28个4KB，需要处理的序号范围是[0, 2^28-1]。每个序号使用1位来表示，则总共需要2^25B即32MB的磁盘空间。
// 2.bitMap的使用场景为较为频繁的随机读写，且以读居多。
// 综合以上两点，使用mmap系统调用来进行文件读写，内存可以完全容下bitmap，因此不会有太多的缺页中断。
type bitMap struct {
	mu sync.RWMutex
	// 负责的最大序号
	maxSeq uint64
	// mmap映射磁盘中的bitmap文件，该文件的每一位用来标识某个序号是否被使用
	mmap util.MMap
}

// todo: 对maxSeq做处理，使其成为8的倍数
func NewBitMap(basePath string, maxSeq uint64) (*bitMap, error) {
	finalPath := util.Join(basePath, "bitmap")
	file, err := os.OpenFile(finalPath, os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	realSeq := uint64(info.Size() * 8)
	if realSeq != 0 && realSeq != maxSeq {
		return nil, fmt.Errorf("exist bitmap max seq: %d, but expect to set %d", realSeq, maxSeq)
	}
	if err := file.Truncate(int64(maxSeq / 8)); err != nil {
		return nil, err
	}
	mmap, err := util.Map(file, util.RDWR, 0)
	if err != nil {
		return nil, err
	}
	return &bitMap{
		maxSeq: maxSeq,
		mmap:   mmap,
	}, nil
}

func (bm *bitMap) IsSet(low, high uint64) bool {
	bm.mu.RLock()
	defer bm.mu.RUnlock()
	for seq := low; seq <= high; seq++ {
		if !bm.isSetOne(seq) {
			return false
		}
	}
	return true
}

func (bm *bitMap) isSetOne(seq uint64) bool {
	ind, offset := seq/8, seq%8
	if (bm.mmap[ind]>>offset)&1 > 0 {
		return true
	}
	return false
}

func (bm *bitMap) Set(low, high uint64) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for seq := low; seq <= high; seq++ {
		bm.setOne(seq)
	}
	return bm.mmap.Flush()
}

func (bm *bitMap) setOne(seq uint64) {
	ind, offset := seq/8, seq%8
	bm.mmap[ind] |= 1 << offset
}

func (bm *bitMap) Remove(low, high uint64) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()
	for seq := low; seq <= high; seq++ {
		bm.removeOne(seq)
	}
	return bm.mmap.Flush()
}

func (bm *bitMap) removeOne(seq uint64) {
	ind, offset := seq/8, seq%8
	bm.mmap[ind] &= ^(1 << offset)
}

type task func()

type chunkFunc func(buf []byte, low, high, chunkIdx uint64) task

type StorageChunk struct {
	bitmap *bitMap
	// default: 4kb
	blockSize uint64
	// todo: 选择一个更合适的大小
	// default: 1gb
	chunkSize     uint64
	blockPerChunk uint64
	// 总chunk数
	chunkNum uint64
	// 总block数
	blockNum uint64
	// key为chunk id
	chunks map[uint64]*Chunk
}

// todo: 对传入的参数做调整
func NewStorageChunk(basePath string, blockSize, chunkSize, blockNum uint64) (storage.Storage, error) {
	bm, err := NewBitMap(basePath, blockNum)
	if err != nil {
		return nil, err
	}
	blockPerChunk := chunkSize / blockSize
	chunkNum := blockNum / blockPerChunk
	chunkMap := make(map[uint64]*Chunk)
	for i := uint64(0); i < chunkNum; i++ {
		newChunk, err := NewChunk(i, basePath, chunkSize)
		if err != nil {
			return nil, err
		}
		chunkMap[i] = newChunk
	}
	return &StorageChunk{
		bitmap:        bm,
		blockSize:     blockSize,
		chunkSize:     chunkSize,
		blockPerChunk: blockPerChunk,
		chunkNum:      chunkNum,
		blockNum:      blockNum,
		chunks:        chunkMap,
	}, nil
}

func (sc *StorageChunk) Get(ctx context.Context, low, high uint64) (*storage.Transfer, error) {
	if low >= sc.blockNum {
		return nil, ErrStorageChunkBeyond
	}
	if !sc.bitmap.IsSet(low, high) {
		return nil, ErrStorageChunkNotExist
	}
	buf := make([]byte, (high-low+1)*sc.blockSize)
	tasks := sc.divideRange(buf, low, high, sc.chunkReadTask)
	for _, tk := range tasks {
		tk()
	}
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
	for _, tk := range tasks {
		tk()
	}
	return sc.bitmap.Set(transfer.Low, transfer.High)
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
