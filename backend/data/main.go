package main

import (
	"flag"
	"fmt"
	"github.com/DMwangnima/easy-disk/data/initialization"
	"net/http"
	"strconv"
)

func main() {
	configPath := flag.String("configPath", "", "path of config file")
	flag.Parse()
	if *configPath == "" {
		panic(fmt.Sprintf("configPath is empty, please run with -configPath <path>"))
	}
	server, err := initialization.InitServer(*configPath)
	if err != nil {
		panic(err)
	}
	if err = http.ListenAndServe(initialization.LocalConfig.Server.ListenIp+":"+strconv.Itoa(initialization.LocalConfig.Server.Port), server); err != nil {
		panic(fmt.Sprintf("server failed, err: %s", err))
	}
}
