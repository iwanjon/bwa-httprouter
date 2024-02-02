package app

import (
	"bwahttprouter/helper"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewDB() *sql.DB {
	db, err := sql.Open("mysql", "root:a@tcp(localhost:5555)/bwastartup?charset=utf8mb4&parseTime=True&loc=Local")

	helper.PanicIfError(err, " error connection db")

	db.SetMaxIdleConns(5)
	db.SetMaxOpenConns(20)
	db.SetConnMaxLifetime(60 * time.Minute)
	db.SetConnMaxIdleTime(10 * time.Minute)
	// defer func() {
	// 	fmt.Println("end the db")
	// 	db.Close()
	// 	fmt.Println("end the db r")
	// }()
	return db

}
