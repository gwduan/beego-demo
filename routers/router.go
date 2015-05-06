package routers

import (
	"beego-demo/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/v1/users/register",
		&controllers.UserController{}, "post:Register")
	beego.Router("/v1/users/login",
		&controllers.UserController{}, "post:Login")
	beego.Router("/v1/users/logout",
		&controllers.UserController{}, "post:Logout")
	beego.Router("/v1/users/passwd",
		&controllers.UserController{}, "post:Passwd")
	beego.Router("/v1/users/uploads",
		&controllers.UserController{}, "post:Uploads")

	beego.Router("/v1/roles/:id",
		&controllers.RoleController{}, "get:GetOne")
	beego.Router("/v1/roles",
		&controllers.RoleController{}, "get:GetAll")
	beego.Router("/v1/roles",
		&controllers.RoleController{}, "post:Post")
	beego.Router("/v1/roles/:id",
		&controllers.RoleController{}, "put:Put")
	beego.Router("/v1/roles/:id",
		&controllers.RoleController{}, "delete:Delete")
}
