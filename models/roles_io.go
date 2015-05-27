package models

type RolePostForm struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
}

type RoleGetOneInfo struct {
	RoleInfo *Role `json:"role"`
}

type RoleGetAllInfo struct {
	RolesInfo []Role `json:"roles"`
}

type RolePutForm struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}
