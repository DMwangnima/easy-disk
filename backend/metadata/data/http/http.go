package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/DMwangnima/easy-disk/metadata/data"
	"io/ioutil"
	"net/http"
)

const (
	httpPrefix = "http://"
)

// todo 后续加入故障转移与重试机制
type DataCenter struct {
	address string
	getPath string
	putPath string
	low     uint64
	high    uint64
}

func (dc *DataCenter) Range() (uint64, uint64) {
	return dc.low, dc.high
}

// todo 考虑是否将http等通信方式抽象成transparnt
func (dc *DataCenter) Get(ctx context.Context, ids ...uint64) ([]*data.File, error) {
	if len(ids) <= 0 {
		return nil, errors.New("ids are empty")
	}
	resp, err := http.Get(dc.getPath)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var getResp GetResp
	if err = json.Unmarshal(buf, &getResp); err != nil {
		return nil, err
	}
	return getResp.Body.Files, nil
}

func (dc *DataCenter) Put(ctx context.Context, files ...*data.File) error {
	if len(files) <= 0 {
		return errors.New("files are empty")
	}
    var putReq PutReq
    putReq.Files = files
    buf, err := json.Marshal(&putReq)
    if err != nil {
    	return err
	}
    resp, err := http.Post(dc.putPath, "application/json", bytes.NewReader(buf))
    if err != nil {
    	return err
	}
	defer resp.Body.Close()
    bodyBuf, err := ioutil.ReadAll(resp.Body)
    if err != nil {
    	return err
	}
	var putResp PutResp
    if err = json.Unmarshal(bodyBuf, &putResp); err != nil {
    	return err
	}
	// 对错误码进行解析
	return nil
}
