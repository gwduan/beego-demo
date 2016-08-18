package models

import (
	"database/sql"
	"fmt"
	"time"

	"beego-demo/models/mymysql"

	"github.com/astaxie/beego"
	"github.com/go-sql-driver/mysql"
)

// Role model definiton.
type Role struct {
	ID       int64     `json:"id,omitempty"`
	Name     string    `json:"name,omitempty"`
	Password string    `json:"password,omitempty"`
	RegDate  time.Time `json:"reg_date,omitempty"`
}

// NewRole alloc and initialize a role.
func NewRole(f *RolePostForm, t time.Time) *Role {
	role := Role{
		ID:       f.ID,
		Name:     f.Name,
		Password: f.Password,
		RegDate:  t}

	return &role
}

// Insert insert a role recode to database.
func (r *Role) Insert() (code int, err error) {
	db := mymysql.Conn()

	st, err := db.Prepare("INSERT INTO roles(id, name, password, reg_date) VALUES(?, ?, ?, ?)")
	if err != nil {
		return ErrDatabase, err
	}
	defer st.Close()

	//if result, err := st.Exec(
	if _, err := st.Exec(r.ID, r.Name, r.Password, r.RegDate); err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			//Duplicate key
			if e.Number == 1062 {
				return ErrDupRows, err
			}

			return ErrDatabase, err
		}

		return ErrDatabase, err
	}

	//r.ID, _ = result.LastInsertId()

	return 0, nil
}

// FindByID query a recode according to input id.
func (r *Role) FindByID(id int64) (code int, err error) {
	db := mymysql.Conn()

	st, err := db.Prepare("SELECT id, name, password, reg_date FROM roles WHERE id = ?")
	if err != nil {
		return ErrDatabase, err
	}
	defer st.Close()

	row := st.QueryRow(id)

	var tmpID sql.NullInt64
	var tmpName sql.NullString
	var tmpPassword sql.NullString
	var tmpRegDate mysql.NullTime
	if err := row.Scan(&tmpID, &tmpName, &tmpPassword, &tmpRegDate); err != nil {
		// Not found.
		if err == sql.ErrNoRows {
			return ErrNotFound, err
		}

		return ErrDatabase, err
	}

	if tmpID.Valid {
		r.ID = tmpID.Int64
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

// ClearPass clear password information.
func (r *Role) ClearPass() {
	r.Password = ""
}

// GetAllRoles query all matched records.
func GetAllRoles(queryVal map[string]string, queryOp map[string]string, order map[string]string, limit int64, offset int64) (records []Role, err error) {
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

	st, err := db.Prepare(sqlStr)
	if err != nil {
		return nil, err
	}
	defer st.Close()

	rows, err := st.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	records = make([]Role, 0, limit)
	for rows.Next() {
		var tmpID sql.NullInt64
		var tmpName sql.NullString
		var tmpPassword sql.NullString
		var tmpRegDate mysql.NullTime
		if err := rows.Scan(&tmpID, &tmpName, &tmpPassword, &tmpRegDate); err != nil {
			return nil, err
		}

		r := Role{}
		if tmpID.Valid {
			r.ID = tmpID.Int64
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

// UpdateByID update a recode accroding to input id.
func (r *Role) UpdateByID(id int64, f *RolePutForm) (code int, err error) {
	db := mymysql.Conn()

	st, err := db.Prepare("UPDATE roles SET name = ?, password = ? WHERE id = ?")
	if err != nil {
		return ErrDatabase, err
	}
	defer st.Close()

	result, err := st.Exec(f.Name, f.Password, id)
	if err != nil {
		return ErrDatabase, err
	}

	num, _ := result.RowsAffected()
	if num > 0 {
		return 0, nil
	}

	return ErrNotFound, nil
}

// DeleteByID delete a record accroding to input id.
func (r *Role) DeleteByID(id int64) (code int, err error) {
	db := mymysql.Conn()

	st, err := db.Prepare("DELETE FROM roles WHERE id = ?")
	if err != nil {
		return ErrDatabase, err
	}
	defer st.Close()

	result, err := st.Exec(id)
	if err != nil {
		return ErrDatabase, err
	}

	num, _ := result.RowsAffected()
	if num > 0 {
		return 0, nil
	}

	return ErrNotFound, nil
}
