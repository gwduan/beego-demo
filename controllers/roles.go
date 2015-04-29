package controllers

import (
	"beego-demo/models"
	"github.com/astaxie/beego"
	"strconv"
)

type RoleController struct {
	beego.Controller
}

func (this *RoleController) GetOne() {
	idStr := this.Ctx.Input.Params[":id"]
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	role := models.Role{}
	if code, err := role.FindById(id); err != nil {
		beego.Debug("FindRoleById:", err)
		if code == 100 {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJson()
		return
	}
	beego.Debug("RoleInfo:", &role)

	this.Data["json"] = &role
	this.ServeJson()
}
