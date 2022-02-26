package local

import (
	"github.com/DMwangnima/easy-disk/data/storage"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

const (
	basePath = "/Users/wangyuxuan/net-disk/data-server-test-dir"
	maxSize = 10000
	evictRatio = 0.5
)


func TestObjectPool_allocateOne(t *testing.T) {
	convey.Convey("Initialize ObjectPool instance", t, func() {
		instance := NewObjectPool(basePath, maxSize, evictRatio)
		convey.Convey("allocateOne allocate new object succeed", func() {
			var testId uint64 = 1
			testFlag := storage.READ
			obj, err := instance.allocateOne(testId, testFlag)
			convey.So(err, assertions.ShouldBeNil)

			realObj := obj.(*Object)
			item, ok := instance.items[realObj.id]
			treeLen := instance.sortTree.Len()
			convey.So(ok, assertions.ShouldBeTrue)
			convey.So(item, assertions.ShouldEqual, realObj)
			convey.So(treeLen, assertions.ShouldBeZeroValue)
			convey.So(item.using, assertions.ShouldEqual, 1)
			convey.So(item.generation, assertions.ShouldEqual, instance.generation)
			convey.So(item.sequence, assertions.ShouldEqual, instance.sequence)
			convey.So(item.flag, assertions.ShouldEqual, testFlag)

			convey.Convey("allocateOne hit cache succeed", func() {
				cacheObj, cacheErr := instance.allocateOne(testId, testFlag)
				convey.So(cacheErr, assertions.ShouldBeNil)
				convey.So(cacheObj, assertions.ShouldEqual, obj)
			})

			convey.Convey("allocateOne flag incompatible", func() {
				cacheObj, cacheErr := instance.allocateOne(testId, storage.WRITE)
				convey.So(cacheObj, assertions.ShouldBeNil)
				convey.So(cacheErr, assertions.ShouldEqual, ErrObjectPoolFlag)
			})
			instance.removeOne(testId)
		})
	})
}

func TestObjectPool_resumeOne(t *testing.T) {
	convey.Convey("Initialize ObjectPool instance", t, func() {
		instance := NewObjectPool(basePath, maxSize, evictRatio)
		convey.Convey("resumeOne succeed without triggering resumeTask", func() {
			var testId uint64 = 1
			testFlag := storage.READ
			obj, err := instance.allocateOne(testId, testFlag)
			convey.So(err, assertions.ShouldBeNil)
			realObj := obj.(*Object)
			convey.So(realObj.using, assertions.ShouldEqual, 1)
			convey.So(realObj.generation, assertions.ShouldEqual, 0)
			convey.So(realObj.sequence, assertions.ShouldEqual, 1)
			instance.resumeOne(obj)
			item, ok := instance.items[realObj.id]
			convey.So(ok, assertions.ShouldBeTrue)
			convey.So(item, assertions.ShouldEqual, realObj)
			convey.So(instance.sortTree.Len(), assertions.ShouldEqual, 1)

			instance.removeOne(realObj.id)
		})
	})
}
