package http

import "github.com/DMwangnima/easy-disk/metadata/data"

type RespBase struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type GetResp struct {
	RespBase
	Body *GetBody `json:"body"`
}

type GetBody struct {
	Files []*data.File `json:"files"`
}

type PutReq struct {
	Files []*data.File `json:"files"`
}

type PutResp struct {
	RespBase
}