package controllers

import (
	"github.com/astaxie/beego"
)

// MainController definition.
type MainController struct {
	beego.Controller
}

// Get method.
func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
	c.Render()
}
