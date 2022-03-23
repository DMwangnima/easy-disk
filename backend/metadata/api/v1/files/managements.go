package files

import (
	"github.com/DMwangnima/easy-disk/metadata/model"
	"github.com/DMwangnima/easy-disk/metadata/model/v1/files"
	"github.com/gin-gonic/gin"
)

// List To list contents of a directory with specified path.
// @Description To list contents of a directory with specified path.
// @Accept      json
// @Produce     json
// @Param path query string true "path of the wanted directory."
// @Param order_by query string true "field by which to order. Options: name, time, size."
// @Param order query string true "sort order. Options: asc, desc."
// @Param limit query int true "number of returns. Not great than 20."
// @Param offset query int true "index of the first returned item."
// @Success 200 {object} files.ListResp
// @Failure 200 {object} model.RespBase
// @Router /api/v1/files/managements/list [get]
func List(ctx *gin.Context) {
	// meaningless sentence. Just let swagger be able to scan the specified struct.
	files.ListResp{}
	model.RespBase{}
}

// Get Download file's partial content by using parent_path, name, and seq(the index of the part of the file).
// @Description Download file's partial content by parent_path, name, and seq(the index of the part of the file). PS: The content-type of response depends on the file category.
// @Accept json
// @Produce application/octet-stream
// @Param parent_path query string true "parent path of the wanted file."
// @Param name query string true "file name."
// @Param seq query int true "the index of file partial. If this file is small than 4MB, just use 0."
// @Success 200
// @Failure 200 {object} model.RespBase
// @Router /api/v1/files/managements/get [get]
func Get(ctx *gin.Context) {

}

// Put Upload a small file(small than 4MB) or a directory directly.
// @Description Upload a small file(small than 4MB) or a directory directly.
// @Accept multipart/form-data
// @Produce json
// @Param parent_path query string true "parent path of the file or directory to upload."
// @Param name query string true "name of the file or directory to upload."
// @Param is_dir query int true "indicate this object is a file or a directory. 0:file, 1:directory."
// @Param hash query int false "md5 of the file. If the upload type is a directory, just miss it."
// @Param file body string false "file content."
// @Success 200 {object} model.RespBase
// @Failure 200 {object} model.RespBase
// @Router /api/v1/files/managements/put [post]
func Put(ctx *gin.Context) {

}

// Delete Delete files or directories asynchronously.
// @Description Delete files or directories asynchronously. Then use token to poll the result.
// @Accept json
// @Produce json
// @Param load body files.DeleteReq true "path list of files or directories to delete."
// @Success 200 {object} files.DeleteResp
// @Failure 200 {object} model.RespBase
// @Router /api/v1/files/managements/delete [post]
func Delete(ctx *gin.Context) {

}

// Rename Rename a file or directory.
// @Description Rename a file or directory.
// @Accept json
// @Produce json
// @Param parent_path query string true "parent path of the file or directory to rename."
// @Param old_name query string true "old name of the file or directory."
// @Param new_name query string true "new name of the file or directory."
// @Success 200 {object} model.RespBase
// @Failure 200 {object} model.RespBase
// @Router /api/v1/files/managements/rename [post]
func Rename(ctx *gin.Context) {

}
