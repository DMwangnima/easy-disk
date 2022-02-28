package model

type GetReq struct {
	Low  uint64 `form:"low"`
	High uint64 `form:"high"`
}

type File struct {
	Low  uint64 `json:"low"`
	High uint64 `json:"high"`
	Data []byte `json:"data"`
}

type GetResp struct {
	File File `json:"file"`
}

type PutReq struct {
	File File `json:"file"`
}

type PutResp struct {
	Error error `json:"codes"`
}
