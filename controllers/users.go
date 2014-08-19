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

func (this *UserController) Login() {
	form := models.LoginForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseLoginForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	beego.Debug("ParseLoginForm:", &form)

	valid := validation.Validation{}
	ok, err := valid.Valid(&form)
	if err != nil {
		beego.Debug("ValidLoginForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	if !ok {
		beego.Debug("ValidLoginForm errors:")
		for _, err := range valid.Errors {
			beego.Debug(err.Key, err.Message)
		}
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	user := models.User{}
	if code, err := user.FindById(form.Phone); err != nil {
		beego.Debug("FindUserById:", err)
		if (code == 100) {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJson()
		return
	}
	beego.Debug("UserInfo:", &user)

	if !user.CheckPass(form.Password) {
		this.Data["json"] = models.NewErrorInfo(ErrPass)
		this.ServeJson()
		return
	}
	user.ClearPass()

	this.SetSession("user_id", form.Phone)

	this.Data["json"] = &models.LoginInfo{Code: 0, UserInfo: &user}
	this.ServeJson()
}
