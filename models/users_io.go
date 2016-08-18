package models

// RegisterForm definiton.
type RegisterForm struct {
	Phone    string `form:"phone"    valid:"Required;Mobile"`
	Name     string `form:"name"     valid:"Required"`
	Password string `form:"password" valid:"Required"`
}

// LoginForm definiton.
type LoginForm struct {
	Phone    string `form:"phone"    valid:"Required;Mobile"`
	Password string `form:"password" valid:"Required"`
}

// LoginInfo definiton.
type LoginInfo struct {
	Code     int   `json:"code"`
	UserInfo *User `json:"user"`
}

// LogoutForm defintion.
type LogoutForm struct {
	Phone string `form:"phone" valid:"Required;Mobile"`
}

// PasswdForm definition.
type PasswdForm struct {
	Phone   string `form:"phone"        valid:"Required;Mobile"`
	OldPass string `form:"old_password" valid:"Required"`
	NewPass string `form:"new_password" valid:"Required"`
}

// UploadsForm definiton.
type UploadsForm struct {
	Phone string `form:"phone" valid:"Required;Mobile"`
}
