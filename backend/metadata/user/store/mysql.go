package store

import (
	"context"
	"database/sql"
	"github.com/DMwangnima/easy-disk/metadata/user"
	"github.com/didi/gendry/builder"
	"github.com/didi/gendry/scanner"
)

const (
	TABLE = "user"
)

var (
    fields []string
)

func init() {
	// todo 反射拿到相应字段，不做硬编码
	fields = []string{"uid", "nick_name", "email", "group", "file_num", "register_date", "last_login_date"}
}

type MysqlManager struct {
    db *sql.DB
}

func (mm *MysqlManager) Get(ctx context.Context, uid string) (*user.User, error) {
	var resUser user.User
    where := map[string]interface{} {
    	"uid": uid,
	}
	sentence, vals, err := builder.BuildSelect(TABLE, where, fields)
	if err != nil {
        return nil, err
	}
	rows, err := mm.db.Query(sentence, vals...)
	if err != nil {
		return nil, err
	}
	if err = scanner.ScanClose(rows, &resUser); err != nil {
		return nil, err
	}
    return &resUser, nil
}

func (mm *MysqlManager) Create(ctx context.Context, opts ...user.UserOption) error {

}

func (mm *MysqlManager) Update(ctx context.Context, uid string, opts ...user.UserOption) error {

}

func (mm *MysqlManager) Delete(ctx context.Context, uid string) error {
    where := map[string]interface{} {
    	"uid": uid,
	}
	sentence, vals, err := builder.BuildDelete(TABLE, where)
	if err != nil {
		return nil
	}
	if _, err = mm.db.Exec(sentence, vals...); err != nil {
		return err
	}
    return nil
}

