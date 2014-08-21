package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/validation"
	"beego-demo/models"
	"time"
	"os"
	"io"
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

func (this *UserController) Logout() {
	form := models.LogoutForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	beego.Debug("ParseLogoutForm:", &form)

	valid := validation.Validation{}
	ok, err := valid.Valid(&form)
	if err != nil {
		beego.Debug("ValidLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	if !ok {
		beego.Debug("ValidLogoutForm errors:")
		for _, err := range valid.Errors {
			beego.Debug(err.Key, err.Message)
		}
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	if this.GetSession("user_id") != form.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJson()
		return
	}

	this.DelSession("user_id")

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJson()
}

func (this *UserController) Passwd() {
	form := models.PasswdForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParsePasswdForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	beego.Debug("ParsePasswdForm:", &form)

	valid := validation.Validation{}
	ok, err := valid.Valid(&form)
	if err != nil {
		beego.Debug("ValidPasswdForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}
	if !ok {
		beego.Debug("ValidPasswdForm errors:")
		for _, err := range valid.Errors {
			beego.Debug(err.Key, err.Message)
		}
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	if this.GetSession("user_id") != form.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJson()
		return
	}

	if code, err := models.ChangePass(form.Phone, form.OldPass,
			form.NewPass); err != nil {
		beego.Debug("ChangeUserPass:", err)
		if (code == 100) {
			this.Data["json"] = models.NewErrorInfo(ErrNoUserPass)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJson()
		return
	}

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJson()
}

func (this *UserController) Uploads() {
	phone := this.GetString("phone")
	beego.Debug("Input phone:", phone)

	if this.GetSession("user_id") != phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJson()
		return
	}

	files := this.Ctx.Request.MultipartForm.File["photos"]
	for i, _ := range files {
		src, err := files[i].Open()
		if err != nil {
			beego.Debug("Open MultipartForm File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrOpenFile)
			this.ServeJson()
			return
		}
		defer src.Close()

		dst, err := os.Create("/var/tmp/" + files[i].Filename)
		if err != nil {
			beego.Debug("Create File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrWriteFile)
			this.ServeJson()
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			beego.Debug("Copy File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrWriteFile)
			this.ServeJson()
			return
		}
	}

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJson()
}
