package models

import (
	"beego-demo/models/mymysql"
	"database/sql"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/go-sql-driver/mysql"
	"time"
)

type Role struct {
	Id       int64     `json:"id"`
	Name     string    `json:"name"`
	Password string    `json:"password"`
	RegDate  time.Time `json:"reg_date"`
}

func NewRole(f *RolePostForm, t time.Time) *Role {
	role := Role{
		Id:       f.Id,
		Name:     f.Name,
		Password: f.Password,
		RegDate:  t}

	return &role
}

func (r *Role) Insert() (code int, err error) {
	db := mymysql.Conn()

	//if result, err := db.Exec(
	if _, err := db.Exec("INSERT INTO roles(id, name, password, reg_date)"+
		" VALUES(?, ?, ?, ?)", r.Id, r.Name, r.Password,
		r.RegDate); err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			//Duplicate key
			if e.Number == 1062 {
				return ErrDupRows, err
			} else {
				return ErrDatabase, err
			}
		} else {
			return ErrDatabase, err
		}
	} else {
		//r.Id, _ = result.LastInsertId()
	}

	return 0, nil
}

func (r *Role) FindById(id int64) (code int, err error) {
	db := mymysql.Conn()

	row := db.QueryRow("SELECT id, name, password, reg_date FROM roles"+
		" WHERE id = ?", id)
	var tmpId sql.NullInt64
	var tmpName sql.NullString
	var tmpPassword sql.NullString
	var tmpRegDate mysql.NullTime
	if err := row.Scan(&tmpId, &tmpName, &tmpPassword,
		&tmpRegDate); err != nil {
		if err == sql.ErrNoRows {
			return ErrNotFound, err
		} else {
			return ErrDatabase, err
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

func (r *Role) ClearPass() {
	r.Password = ""
}

func GetAllRoles(queryVal map[string]string, queryOp map[string]string,
	order map[string]string, limit int64,
	offset int64) (records []Role, err error) {
	sqlStr := "SELECT id, name, password, reg_date FROM roles"
	if len(queryVal) > 0 {
		sqlStr += " WHERE "
		first := true
		for k, v := range queryVal {
			if !first {
				sqlStr += " AND "
			} else {
				first = false
			}

			sqlStr += k
			sqlStr += " "
			sqlStr += queryOp[k]
			sqlStr += " '"
			sqlStr += v
			sqlStr += "'"
		}
	}
	if len(order) > 0 {
		sqlStr += " ORDER BY "
		first := true
		for k, v := range order {
			if !first {
				sqlStr += ", "
			} else {
				first = false
			}

			sqlStr += k
			sqlStr += " "
			sqlStr += v
		}
	}
	sqlStr += " LIMIT " + fmt.Sprintf("%d", limit)
	if offset > 0 {
		sqlStr += " OFFSET " + fmt.Sprintf("%d", offset)
	}
	beego.Debug("sqlStr:", sqlStr)

	db := mymysql.Conn()
	rows, err := db.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records = make([]Role, 0, limit)
	for rows.Next() {
		var tmpId sql.NullInt64
		var tmpName sql.NullString
		var tmpPassword sql.NullString
		var tmpRegDate mysql.NullTime
		if err := rows.Scan(&tmpId, &tmpName, &tmpPassword,
			&tmpRegDate); err != nil {
			return nil, err
		}

		r := Role{}
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
		records = append(records, r)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return records, nil
}

func (r *Role) UpdateById(id int64, f *RolePutForm) (code int, err error) {
	db := mymysql.Conn()

	result, err := db.Exec("UPDATE roles SET name = ?, password = ?"+
		" WHERE id = ?", f.Name, f.Password, id)
	if err != nil {
		return ErrDatabase, err
	}

	num, _ := result.RowsAffected()
	if num > 0 {
		return 0, nil
	} else {
		return ErrNotFound, nil
	}
}

func (r *Role) DeleteById(id int64) (code int, err error) {
	db := mymysql.Conn()

	result, err := db.Exec("DELETE FROM roles WHERE id = ?", id)
	if err != nil {
		return ErrDatabase, err
	}

	num, _ := result.RowsAffected()
	if num > 0 {
		return 0, nil
	} else {
		return ErrNotFound, nil
	}
}
