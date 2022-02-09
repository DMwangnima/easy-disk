package local

import "github.com/DMwangnima/easy-disk/data/storage"

type Stream struct {
	channel chan *storage.Transfer
	err     error
}

func NewStream(chanSize int) storage.Stream {
	return &Stream{
		channel: make(chan *storage.Transfer, chanSize),
	}
}

func (s *Stream) Consume() (*storage.Transfer, bool) {
    trans, ok := <-s.channel
    return trans, ok
}

// 考虑是否需要异步
func (s *Stream) Produce(trans *storage.Transfer) {
    s.channel <- trans
}

func (s *Stream) Error() string {
	return s.err.Error()
}
