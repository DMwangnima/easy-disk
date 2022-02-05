package model

type GetReq struct {
	Low  uint64 `form:"low"`
	High uint64 `form:"high"`
}

type File struct {
	Id   uint64 `json:"id"`
	Data []byte `json:"data"`
}

type GetResp struct {
	Files []File `json:"files"`
}

type PutReq struct {
	Files []File `json:"files"`
}

type PutResp struct {
	Error error `json:"code"`
}