package v1

import (
    "github.com/DMwangnima/easy-disk/data/api"
    "github.com/DMwangnima/easy-disk/data/codes"
    "github.com/DMwangnima/easy-disk/data/model"
    "github.com/DMwangnima/easy-disk/data/storage"
    "github.com/DMwangnima/easy-disk/data/util"
    "github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
    var getReq model.GetReq
    if err := ctx.ShouldBindQuery(&getReq); err != nil {
        ctx.JSON(400, model.FailureWithCode(codes.WRONG_PARAM))
        return
    }
    sli := util.GenerateContinuousSlice(getReq.Low, getReq.High)
    files, err := api.Store.Get(ctx, sli...)
    if err != nil {
        // 统一返回错误并打日志
    }
    resp := generateGetResp(files)
    ctx.JSON(200, model.SuccessWithBody(resp))
}

func generateGetResp(transfers []*storage.Transfer) model.GetResp {
    files := make([]model.File, len(transfers))
    for i := 0; i < len(transfers); i++ {
        files[i] = model.File{
            Id:   transfers[i].Id,
            Data: transfers[i].Data,
        }
    }
    return model.GetResp{
        Files:    files,
    }
}

func Put(ctx *gin.Context) {
    var putReq model.PutReq
    if err := ctx.ShouldBindJSON(&putReq); err != nil {
        ctx.JSON(400, model.FailureWithCode(codes.WRONG_PARAM))
        return
    }
    transfers := generateTransfers(&putReq)
    if err := api.Store.Put(ctx, transfers...); err != nil {

    }
    ctx.JSON(201, model.Success())
}

func generateTransfers(putReq *model.PutReq) []*storage.Transfer {
    res := make([]*storage.Transfer, len(putReq.Files))
    for i, file := range putReq.Files {
        res[i] = &storage.Transfer{
            Id:   file.Id,
            Data: file.Data,
        }
    }
    return res
}
