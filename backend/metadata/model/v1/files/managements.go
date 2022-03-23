package files

import (
	"github.com/DMwangnima/easy-disk/metadata/model"
	"time"
)

type ListReq struct {
	Path    string `param:"path"`
	OrderBy string `param:"order_by"`
	Order   string `param:"order"`
	Limit   int    `param:"limit"`
	Offset  int    `param:"offset"`
}

type File struct {
	ParentPath string    `json:"parent_path"`
	Name       string    `json:"name"`
	Size       int       `json:"size"`
	CreateTime time.Time `json:"create_time"`
	UpdateTime time.Time `json:"update_time"`
	IsDir      int       `json:"is_dir"`
	Hash       string    `json:"hash"`
	Category   int       `json:"category"`
}

type ListResp struct {
	model.RespBase
	Files []*File `json:"files"`
}

type DeleteReq struct {
	FileList []string `json:"file_list"`
}

type DeleteResp struct {
	model.RespBase
	TaskToken string `json:"task_token"`
}
