package user

import (
	"bwahttprouter/exception"
	"bwahttprouter/helper"
	"context"
	"database/sql"
	"fmt"
)

type Repository interface {
	Save(ctx context.Context, user User) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	FindById(ctx context.Context, id int) (User, error)
	Update(ctx context.Context, user User) (User, error)
}

type repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *repository {

	return &repository{db}
}

func (r *repository) Save(ctx context.Context, user User) (User, error) {
	// var w http.ResponseWriter
	tx, err := r.db.Begin()
	helper.PanicIfError(err, " error repo save tx")
	SQL := "insert into users(Name, Occupation, Email, Password_Hash, Avatar_File_Name, Role) values (?,?,?,?,?,?)"
	result, err := tx.ExecContext(ctx, SQL, user.Name, user.Occupation, user.Email, user.PasswordHash, user.AvatarFileName, user.Role)
	// SQL := "insert into users(Name, Occupation, Email, PasswordHash, AvatarFileName, Role, CreatedAt, UpdatedAt) values (?,?,?,?,?,?,?,?)"
	// result, err := tx.ExecContext(ctx, SQL, user.Name, user.Occupation, user.Email, user.PasswordHash, user.AvatarFileName, user.Role, user.CreatedAt, user.UpdatedAt)

	helper.PanicIfError(err, "error repo save result")
	defer helper.CommitOrRollback(tx)

	id, err := result.LastInsertId()
	helper.PanicIfError(err, "error repo save id")

	user.ID = int(id)

	return user, nil
}

func (r *repository) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error repo find email tx")
	SQL := "select * from users where email = ?"
	result := tx.QueryRowContext(ctx, SQL, email)
	err = result.Err()
	helper.PanicIfError(err, "error repo find email reult")
	defer helper.CommitOrRollback(tx)
	err = result.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	fmt.Println(err == sql.ErrNoRows, "eeeee")
	exception.PanicIfNotFound(err, "errr repo scan find by email")
	// helper.PanicIfError(err, "error repo find email scan")
	return user, nil
}

func (r *repository) FindById(ctx context.Context, id int) (User, error) {
	var user User
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error find y id repository")
	SQL := "select * from users where id = ? limit 1;"
	result := tx.QueryRowContext(ctx, SQL, id)
	err = result.Err()
	helper.PanicIfError(err, "error find by id repository")
	defer helper.CommitOrRollback(tx)
	err = result.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	// helper.PanicIfError(err, "error  find email scan")
	exception.PanicIfNotFound(err, "errr repo scan find by id")
	return user, nil

}
func (r *repository) Update(ctx context.Context, user User) (User, error) {
	tx, err := r.db.Begin()
	helper.PanicIfError(err, "error find y id repository")
	// user.ID = 333
	// user.Role = 11
	SQL := "update users set name = ? , occupation = ? , email = ? , password_hash = ? , avatar_file_name = ? , role = ?, token = ? where id = ?"
	// SQL := "select * from users where id = ? limit 1;"
	_, err = tx.ExecContext(ctx, SQL, user.Name, user.Occupation, user.Email, user.PasswordHash, user.AvatarFileName, user.Role, user.Token, user.ID)
	fmt.Println(err, "madang")
	// err = result.Err()
	// helper.PanicIfError(err, "error find by id repository")
	exception.PanicIfNotFound(err, "errr repo scan update")
	defer helper.CommitOrRollback(tx)
	// err = result.Scan(&user.ID, &user.Name, &user.Occupation, &user.Email, &user.PasswordHash, &user.AvatarFileName, &user.Role, &user.Token, &user.CreatedAt, &user.UpdatedAt)
	// helper.PanicIfError(err, "error  find email scan")

	return user, nil
}
