package files

import "github.com/DMwangnima/easy-disk/metadata/model"

type StartResp struct {
	model.RespBase
	SessionToken string `json:"session_token"`
	BlockList []int `json:"block_list" example:"0,1,2,3"`
}

type FinishReq struct {
	BlockList []string `json:"block_list"`
}