package model

import "github.com/DMwangnima/easy-disk/data/codes"

type RespBase struct {
	Message string `json:"message"`
	Code    int    `json:"codes"`
}

type Response struct {
	RespBase
	Body interface{} `json:"body"`
}

func Success() Response {
	return Response{
		RespBase: RespBase{
			Message: codes.CodeMap[codes.OK],
			Code: codes.OK,
		},
	}
}

func SuccessWithBody(body interface{}) Response {
	return Response{
		RespBase: RespBase{
			Message: codes.CodeMap[codes.OK],
			Code:    codes.OK,
		},
		Body: body,
	}
}

func FailureWithDetail(code int, body interface{}) Response {
	code, msg := codes.GetMessage(code)
	return Response{
		RespBase: RespBase{
			Message: msg,
			Code:    code,
		},
		Body: body,
	}
}

func FailureWithCode(code int) Response {
	code, msg := codes.GetMessage(code)
	return Response{
		RespBase: RespBase{
			Message: msg,
			Code: code,
		},
	}
}
