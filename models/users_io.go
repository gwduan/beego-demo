package models

type RegisterForm struct {
	Phone        string `form:"phone"         valid:"Required;Mobile"`
	Name         string `form:"name"          valid:"Required"`
	Password     string `form:"password"      valid:"Required"`
}

type LoginForm struct {
	Phone        string `form:"phone"         valid:"Required;Mobile"`
	Password     string `form:"password"      valid:"Required"`
}

type LoginInfo struct {
	Code     int    `json:"code"`
	UserInfo *User  `json:"user"`
}
