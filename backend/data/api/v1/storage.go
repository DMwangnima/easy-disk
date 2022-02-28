package v1

import (
    "github.com/DMwangnima/easy-disk/data/api"
    "github.com/DMwangnima/easy-disk/data/codes"
    "github.com/DMwangnima/easy-disk/data/model"
    "github.com/DMwangnima/easy-disk/data/storage"
    "github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
    var getReq model.GetReq
    if err := ctx.ShouldBindQuery(&getReq); err != nil {
        ctx.JSON(400, model.FailureWithCode(codes.WRONG_PARAM))
        return
    }
    file, err := api.Store.Get(ctx, getReq.Low, getReq.High)
    if err != nil {
        ctx.JSON(500, model.FailureWithCode(codes.STORAGE))
        return
    }
    ctx.JSON(200, model.SuccessWithBody(model.GetResp{
        File: model.File{
            Low:  file.Low,
            High: file.High,
            Data: file.Data,
        },
    }))
}

func Put(ctx *gin.Context) {
    var putReq model.PutReq
    if err := ctx.ShouldBindJSON(&putReq); err != nil {
        ctx.JSON(400, model.FailureWithCode(codes.WRONG_PARAM))
        return
    }
    if len(putReq.File.Data) <= 0 {
        ctx.JSON(400, model.FailureWithCode(codes.WRONG_PARAM))
        return
    }
    if err := api.Store.Put(ctx, &storage.Transfer{
        Low:  putReq.File.Low,
        High: putReq.File.High,
        Data: putReq.File.Data,
    }); err != nil {
        // todo: 统一返回错误
        return
    }
    ctx.JSON(201, model.Success())
}