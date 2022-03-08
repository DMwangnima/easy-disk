package util

import (
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestFloorPowerOf2(t *testing.T) {
	convey.Convey("test FloorPowerOf2", t, func() {
		tests := []struct {
			input  uint64
			expect uint64
		}{
			{input: 1, expect: 1},
			{input: 3, expect: 2},
			{input: 127, expect: 64},
			{input: 1<<63 + 1, expect: 1 << 63},
		}
		for _ ,test := range tests {
			convey.So(FloorPowerOf2(test.input), assertions.ShouldEqual, test.expect)
		}
	})
}
