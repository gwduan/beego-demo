package models

type RegisterForm struct {
	Phone    string `form:"phone"    valid:"Required;Mobile"`
	Name     string `form:"name"     valid:"Required"`
	Password string `form:"password" valid:"Required"`
}

type LoginForm struct {
	Phone    string `form:"phone"    valid:"Required;Mobile"`
	Password string `form:"password" valid:"Required"`
}

type LoginInfo struct {
	Code     int   `json:"code"`
	UserInfo *User `json:"user"`
}

type LogoutForm struct {
	Phone string `form:"phone" valid:"Required;Mobile"`
}

type PasswdForm struct {
	Phone   string `form:"phone"        valid:"Required;Mobile"`
	OldPass string `form:"old_password" valid:"Required"`
	NewPass string `form:"new_password" valid:"Required"`
}

type UploadsForm struct {
	Phone string `form:"phone" valid:"Required;Mobile"`
}
