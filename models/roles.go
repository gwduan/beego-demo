package models

import (
	"beego-demo/models/mymysql"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"time"
)

type Role struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	RegDate  time.Time `json:"reg_date"`
}

func (r *Role) FindById(id int64) (code int, err error) {
	db := mymysql.Conn()

	row := db.QueryRow(
		"SELECT id, name, password, reg_date FROM roles WHERE id = ?",
		id)

	var tmpId sql.NullInt64
	var tmpName sql.NullString
	var tmpPassword sql.NullString
	var tmpRegDate mysql.NullTime
	if err := row.Scan(&tmpId, &tmpName, &tmpPassword,
		&tmpRegDate); err != nil {
		if err == sql.ErrNoRows {
			return 100, err
		} else {
			return -1, err
		}
	}

	if tmpId.Valid {
		r.Id = tmpId.Int64
	}
	if tmpName.Valid {
		r.Name = tmpName.String
	}
	if tmpPassword.Valid {
		r.Password = tmpPassword.String
	}
	if tmpRegDate.Valid {
		r.RegDate = tmpRegDate.Time
	}

	return 0, nil
}
