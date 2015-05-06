package models

type RolePostForm struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RoleGetOneInfo struct {
	Code     int   `json:"code"`
	RoleInfo *Role `json:"role"`
}

type RoleGetAllInfo struct {
	Code      int    `json:"code"`
	RolesInfo []Role `json:"roles"`
}

type RolePutForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
