package model

import "github.com/DMwangnima/easy-disk/data/code"

type RespBase struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

type Response struct {
	RespBase
	Body interface{} `json:"body"`
}

func SuccessWithBody(body interface{}) Response {
    return Response{
        RespBase: RespBase{
        	Message: code.CodeMap[code.OK],
        	Code: code.OK,
        },
		Body: body,
	}
}
