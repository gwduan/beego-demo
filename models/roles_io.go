package models

// RolePostForm definiton.
type RolePostForm struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// RolePostInfo definiton.
type RolePostInfo struct {
	RoleInfo *Role `json:"role"`
}

// RoleGetOneInfo definiton.
type RoleGetOneInfo struct {
	RoleInfo *Role `json:"role"`
}

// RoleGetAllInfo definiton.
type RoleGetAllInfo struct {
	RolesInfo []Role `json:"roles"`
}

// RolePutForm definiton.
type RolePutForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// RolePutInfo definiton.
type RolePutInfo struct {
	RoleInfo *Role `json:"role"`
}

// RoleAuthForm definiton.
type RoleAuthForm struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

// RoleAuthInfo definiton.
type RoleAuthInfo struct {
	Token string `json:"token"`
}
