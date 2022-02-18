package local

import "github.com/DMwangnima/easy-disk/data/storage"

type Stream struct {
	channel chan storage.Object
	err     error
}

func NewStream(chanSize int) storage.Stream {
	return &Stream{
		channel: make(chan storage.Object, chanSize),
	}
}

func (s *Stream) Consume() (storage.Object, bool) {
    obj, ok := <-s.channel
    return obj, ok
}

// 考虑是否需要异步
func (s *Stream) Produce(obj storage.Object) {
    s.channel <- obj
}

func (s *Stream) Error() string {
	return s.err.Error()
}
