package models

type RegisterForm struct {
	Phone        string `form:"phone"         valid:"Required;Mobile"`
	Name         string `form:"name"          valid:"Required"`
	Password     string `form:"password"      valid:"Required"`
}
