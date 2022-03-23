package files

import (
	"github.com/DMwangnima/easy-disk/metadata/model"
	"github.com/DMwangnima/easy-disk/metadata/model/v1/files"
	"github.com/gin-gonic/gin"
)

// Start Start a session to upload a big file(great than 4MB).
// @Description Start a session to upload a big file(great than 4MB).
// @Tags sessions
// @Accept json
// @Produce json
// @Param parent_path query string true "parent path of big file to upload."
// @Param name query string true "name of big file."
// @Param upload_type query int true "additional effect of this upload action."
// @Param seqs query int true "number of file blocks.(every block is 4MB size)"
// @Success 200 {object} files.StartResp
// @Router /api/v1/files/sessions/start [get]
func Start(ctx *gin.Context) {
	// meaningless sentence. Just let swagger be able to scan the specified struct.
	files.StartResp{}
	model.RespBase{}
}

// Upload Upload a 4MB(or less than) block of the big file in the specified session.
// @Description Upload a 4MB(or less than) block of the big file in the specified session.
// @Tags sessions
// @Accept multipart/form-data
// @Produce json
// @Param session_token query string true "token of the session."
// @Param seq query int true "index of the block. eg: the first block of the big file has the index 0."
// @Param hash query string true "md5 of this block."
// @Param file body string true "file content"
// @Success 200 {object} model.RespBase
// @Router /api/v1/files/sessions/upload [post]
func Upload(ctx *gin.Context) {

}

// FastUpload Send the md5 of the whole file to check whether this file has been uploaded.
// @Description Send the md5 of the whole file to check whether this file has been uploaded.
// @Tags sessions
// @Accept json
// @Produce json
// @Param session_token query string true "token of the session."
// @Param hash query string true "md5 of the whole file."
// @Success 200 {object} model.RespBase
// @Router /api/v1/files/sessions/fastUpload [get]
func FastUpload(ctx *gin.Context) {

}

// Finish Finish the specified session and create the big file in fileSystem.
// @Description Finish the specified session and create the big file in fileSystem.
// @Tags sessions
// @Accept json
// @Produce json
// @Param session_token query string true "token of the session."
// @Param load body files.FinishReq true "block hash list belongs to the whole big file. eg: ["312321", "123124"] this file consists of two blocks."
// @Success 200 {object} model.RespBase
// @Router /api/v1/files/sessions/finish [post]
func Finish(ctx *gin.Context) {

}
