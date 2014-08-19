package models

import (
	"beego-demo/models/mymongo"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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

func (u *User) FindById(id string) (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB("").C("users")
	err = c.FindId(id).One(u);

	if err != nil {
		if err == mgo.ErrNotFound {
			code = 100
		} else {
			code = -1
		}
	} else {
		code = 0
	}
	return
}

func (u *User) CheckPass(pass string) bool {
	return u.Password == pass
}

func (u *User) ClearPass() {
	u.Password = ""
}

func ChangePass(id, oldPass, newPass string) (code int, err error) {
	mConn := mymongo.Conn()
	defer mConn.Close()

	c := mConn.DB("").C("users")
	err = c.Update(bson.M{"_id": id, "password": oldPass},
			bson.M{"$set": bson.M{"password": newPass}})

	if err != nil {
		if err == mgo.ErrNotFound {
			code = 100
		} else {
			code = -1
		}
	} else {
		code = 0
	}
	return
}
