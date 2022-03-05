package local

import (
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestBitMap(t *testing.T) {
	convey.Convey("initialize bitMap instance", t, func() {
		bm, err := NewBitMap(ChunkPath, 1024)
		if err != nil {
			t.Fatal(err)
		}
		convey.Convey("Set", func() {
			convey.So(bm.Set(0, 7), assertions.ShouldBeNil)
		})
		convey.Convey("IsSet", func() {
			convey.So(bm.IsSet(0, 7), assertions.ShouldBeTrue)
			convey.So(bm.IsSet(8, 9), assertions.ShouldBeFalse)
		})
		convey.Convey("Remove", func() {
			convey.So(bm.Remove(0, 7), assertions.ShouldBeNil)
			convey.So(bm.IsSet(0, 7), assertions.ShouldBeFalse)
		})
	})
}
