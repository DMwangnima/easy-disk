package util

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"runtime"
)

//判断文件或文件夹是否存在
func isExist(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		if os.IsNotExist(err) {
			return false
		}
		fmt.Println(err)
		return false
	}
	return true
}

func Join(basePath string, endPath string) string {
	if ok := isExist(basePath); !ok {
		os.MkdirAll(basePath, 0666)
	}
	newPath := path.Join(basePath, endPath)
	if runtime.GOOS == "windows" {
		return filepath.FromSlash(newPath)
	}
	return newPath
}
