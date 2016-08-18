package routers

import (
	"beego-demo/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/", &controllers.MainController{})
	ns := beego.NewNamespace("/v1",
		beego.NSNamespace("/users",
			beego.NSRouter("/register", &controllers.UserController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.UserController{}, "post:Login"),
			beego.NSRouter("/logout", &controllers.UserController{}, "post:Logout"),
			beego.NSRouter("/passwd", &controllers.UserController{}, "post:Passwd"),
			beego.NSRouter("/uploads", &controllers.UserController{}, "post:Uploads"),
			beego.NSRouter("/downloads", &controllers.UserController{}, "get:Downloads"),
		),
		beego.NSNamespace("/roles",
			beego.NSRouter("/:id", &controllers.RoleController{}, "get:GetOne;put:Put;delete:Delete"),
			beego.NSRouter("/", &controllers.RoleController{}, "get:GetAll;post:Post"),
			beego.NSRouter("/auth", &controllers.RoleController{}, "post:Auth"),
		),
	)
	beego.AddNamespace(ns)
}
