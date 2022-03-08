package initialization

import (
	"github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestInitDataConfig(t *testing.T) {
	configPath := "../config/config.yaml"
	initDataConfig(configPath)
	convey.Convey("test", t, func() {
		convey.So(LocalConfig.Storage.BasePath, convey.ShouldEqual, "/Users/wangyuxuan/net-disk/data-server-test-dir/chunk")
		convey.So(LocalConfig.Server.ListenIp, convey.ShouldEqual, "0.0.0.0")
		convey.So(LocalConfig.Server.Port, convey.ShouldEqual, 9000)
		convey.So(LocalConfig.Server.LogPath, convey.ShouldEqual, "/Users")
	})
}
