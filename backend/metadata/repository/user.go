package repository

import (
	"github.com/DMwangnima/easy-disk/metadata/model"
	"github.com/jmoiron/sqlx"
)

type Users struct {
	db *sqlx.DB
	table string
	getByUidStmt *sqlx.Stmt
	insertStmt *sqlx.NamedStmt
	deleteStmt *sqlx.Stmt
}

func NewUsers() *Users {

}

func (users *Users) initGetByUid() {
	var err error
	users.getByUidStmt, err = users.db.Preparex(
		"SELECT id, nick_name, email, user_group, file_num, register_date, last_login_date " +
		"FROM " + users.table + " " +
		"WHERE id = ? " +
		"LIMIT 1")
	if err != nil {
		panic(err)
	}
}

func (users *Users) GetByUid(id string) (*model.User, error) {
	var resUser model.User
	var err error
	if err = users.db.Get(&resUser, id); err != nil {
		return nil, err
	}
	return &resUser, nil
}

func (users *Users) initInsert() {
	var err error
	users.insertStmt, err = users.db.PrepareNamed(
		"INSERT INTO " + users.table + " " +
		"(id, nick_name, email, encrypted_password, user_group, file_num, register_date, last_login_date) " +
		"VALUES (:id, :nick_name, :email, :encrypted_password, :user_group, :file_num, :register_date, :last_login_date)",
		)
	if err != nil {
		panic(err)
	}
}

func (users *Users) Insert(user *model.User) error {
    if _, err := users.insertStmt.Exec(user); err != nil {
    	return err
	}
	return nil
}

func (users *Users) initUpdate() {

}

func (users *Users) Update(user *model.User) error {

}


func (users *Users) initDelete() {
    var err error
    users.deleteStmt, err = users.db.Preparex(
    	"DELETE FROM " + users.table + "" +
    	"WHERE id = ?",
    	)
    if err != nil {
    	panic(err)
	}
}

func (users *Users) Delete(id string) error {
    if _, err := users.deleteStmt.Exec(id); err != nil {
    	return err
	}
	return nil
}