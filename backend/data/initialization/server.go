package initialization

import (
	"github.com/DMwangnima/easy-disk/data/router"
	"net/http"
)

func InitServer(configPath string) (http.Handler, error) {
	var err error
	if err = initDataConfig(configPath); err != nil {
		return nil, err
	}
	if err = initInject(); err != nil {
		return nil ,err
	}
	engine := router.InitRouters()
	return engine, nil
}
