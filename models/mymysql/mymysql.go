package mymysql

import (
	"database/sql"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

var session *sql.DB

func Conn() *sql.DB {
	return session
}

func init() {
	url := beego.AppConfig.String("mysql::url")

	db, err := sql.Open("mysql", url)
	if err != nil {
		panic(err)
	}

	if err := db.Ping(); err != nil {
		panic(err)
	}

	session = db
}
