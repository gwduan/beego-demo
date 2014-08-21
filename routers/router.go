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
}
