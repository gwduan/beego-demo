package models

import (
	"beego-demo/models/mymongo"
	"gopkg.in/mgo.v2"
	"time"
)

type User struct {
	ID           string    `bson:"_id"           json:"_id"`
	Name         string    `bson:"name"          json:"name"`
	Password     string    `bson:"password"      json:"password"`
	RegDate      time.Time `bson:"reg_date"      json:"reg_date"`
}

func NewUser(r *RegisterForm, t time.Time) *User {
	user := User{
		ID:           r.Phone,
		Name:         r.Name,
		Password:     r.Password}
	user.RegDate = t

	return &user
}

func (u *User) Insert() (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB("").C("users")
	err = c.Insert(u)

	if err != nil {
		if  mgo.IsDup(err) {
			code = 100
		} else {
			code = -1
		}
	} else {
		code = 0
	}
	return
}
