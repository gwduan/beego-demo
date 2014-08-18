package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"beego-demo/models"
	"time"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Register() {
	form := models.RegisterForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	beego.Debug("ParseRegsiterForm:", &form)

	valid := validation.Validation{}
	ok, err := valid.Valid(&form)
	if err != nil {
		beego.Debug("ValidRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	if !ok {
		beego.Debug("ValidRegsiterForm errors:")
		for _, err := range valid.Errors {
			beego.Debug(err.Key, err.Message)
		}
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	regDate := time.Now()
	user := models.NewUser(&form, regDate)
	beego.Debug("NewUser:", user)

	if code, err := user.Insert(); err != nil {
		beego.Debug("InsertUser:", err)
		if (code == 100) {
			this.Data["json"] = models.NewErrorInfo(ErrDupUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJson()
		return
	}

	go models.IncTotalUserCount(regDate)

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJson()
}
