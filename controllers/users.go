package controllers

import (
	"beego-demo/models"
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type UserController struct {
	BaseController
}

func (this *UserController) Register() {
	form := models.RegisterForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseRegsiterForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidRegsiterForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	regDate := time.Now()
	user, err := models.NewUser(&form, regDate)
	if err != nil {
		beego.Error("NewUser:", err)
		this.Data["json"] = models.NewErrorInfo(ErrSystem)
		this.ServeJSON()
		return
	}
	beego.Debug("NewUser:", user)

	if code, err := user.Insert(); err != nil {
		beego.Error("InsertUser:", err)
		if code == models.ErrDupRows {
			this.Data["json"] = models.NewErrorInfo(ErrDupUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
	}

	go models.IncTotalUserCount(regDate)

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJSON()
}

func (this *UserController) Login() {
	form := models.LoginForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseLoginForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseLoginForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidLoginForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	user := models.User{}
	if code, err := user.FindById(form.Phone); err != nil {
		beego.Error("FindUserById:", err)
		if code == models.ErrNotFound {
			this.Data["json"] = models.NewErrorInfo(ErrNoUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJSON()
		return
	}
	beego.Debug("UserInfo:", &user)

	if ok, err := user.CheckPass(form.Password); err != nil {
		beego.Error("CheckUserPass:", err)
		this.Data["json"] = models.NewErrorInfo(ErrSystem)
		this.ServeJSON()
		return
	} else if !ok {
		this.Data["json"] = models.NewErrorInfo(ErrPass)
		this.ServeJSON()
		return
	}
	user.ClearPass()

	this.SetSession("user_id", form.Phone)

	this.Data["json"] = &models.LoginInfo{Code: 0, UserInfo: &user}
	this.ServeJSON()
}

func (this *UserController) Logout() {
	form := models.LogoutForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseLogoutForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidLogoutForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	if this.GetSession("user_id") != form.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	this.DelSession("user_id")

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJSON()
}

func (this *UserController) Passwd() {
	form := models.PasswdForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParsePasswdForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParsePasswdForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidPasswdForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	if this.GetSession("user_id") != form.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	code, err := models.ChangePass(form.Phone, form.OldPass, form.NewPass)
	if err != nil {
		beego.Error("ChangeUserPass:", err)
		if code == models.ErrNotFound {
			this.Data["json"] = models.NewErrorInfo(ErrNoUserPass)
		} else if code == models.ErrDatabase {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrSystem)
		}
		this.ServeJSON()
		return
	}

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJSON()
}

func (this *UserController) Uploads() {
	form := models.UploadsForm{}
	if err := this.ParseForm(&form); err != nil {
		beego.Debug("ParseUploadsForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	beego.Debug("ParseUploadsForm:", &form)

	if err := this.VerifyForm(&form); err != nil {
		beego.Debug("ValidUploadsForm:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}

	if this.GetSession("user_id") != form.Phone {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	//files := this.Ctx.Request.MultipartForm.File["photos"]
	files, err := this.GetFiles("photos")
	if err != nil {
		beego.Debug("GetFiles:", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJSON()
		return
	}
	for i, _ := range files {
		src, err := files[i].Open()
		if err != nil {
			beego.Error("Open MultipartForm File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrOpenFile)
			this.ServeJSON()
			return
		}
		defer src.Close()

		hash := md5.New()
		if _, err := io.Copy(hash, src); err != nil {
			beego.Error("Copy File to Hash:", err)
			this.Data["json"] = models.NewErrorInfo(ErrWriteFile)
			this.ServeJSON()
			return
		}
		hex := fmt.Sprintf("%x", hash.Sum(nil))

		dst, err := os.Create(beego.AppConfig.String("apppath") +
			"static/" + hex + filepath.Ext(files[i].Filename))
		if err != nil {
			beego.Error("Create File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrWriteFile)
			this.ServeJSON()
		}
		defer dst.Close()

		src.Seek(0, 0)
		if _, err := io.Copy(dst, src); err != nil {
			beego.Error("Copy File:", err)
			this.Data["json"] = models.NewErrorInfo(ErrWriteFile)
			this.ServeJSON()
			return
		}
	}

	this.Data["json"] = models.NewNormalInfo("Succes")
	this.ServeJSON()
}

func (this *UserController) Downloads() {
	if this.GetSession("user_id") == nil {
		this.Data["json"] = models.NewErrorInfo(ErrInvalidUser)
		this.ServeJSON()
		return
	}

	file := beego.AppConfig.String("apppath") + "logs/test.log"
	http.ServeFile(this.Ctx.ResponseWriter, this.Ctx.Request, file)
}
