package controllers

import (
	"beego-demo/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	"strconv"
	"time"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Auth() {
	form := models.RoleAuthForm{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRoleAuth:", err)
		this.RetError(errInputData)
		return
	}
	beego.Debug("ParseRoleAuth:", &form)

	role := models.Role{}
	if code, err := role.FindById(form.Id); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			this.RetError(errNoUser)
		} else {
			this.RetError(errDatabase)
		}
		return
	}
	beego.Debug("RoleInfo:", &role)

	if role.Name != form.Name || role.Password != form.Password {
		this.RetError(errPass)
		return
	}

	// Create the token with some claims
	claims := make(jwt.MapClaims)
	claims["id"] = strconv.FormatInt(form.Id, 10)
	claims["name"] = form.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		beego.Error("jwt.SignedString:", err)
		this.RetError(errSystem)
		return
	}

	this.Data["json"] = &models.RoleAuthInfo{Token: tokenString}
	this.ServeJSON()
}

func (this *RoleController) Post() {
	token, e := this.ParseToken()
	if e != nil {
		this.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		if claims["id"] != "1" {
			this.RetError(errPermission)
			return
		}
	} else {
		this.RetError(errPermission)
		return
	}

	form := models.RolePostForm{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePost:", err)
		this.RetError(errInputData)
		return
	}
	beego.Debug("ParseRolePost:", &form)

	regDate := time.Now()
	role := models.NewRole(&form, regDate)
	beego.Debug("NewRole:", role)

	if code, err := role.Insert(); err != nil {
		beego.Error("InsertRole:", err)
		if code == models.ErrDupRows {
			this.RetError(errDupUser)
		} else {
			this.RetError(errDatabase)
		}
		return
	}

	role.ClearPass()

	this.Data["json"] = &models.RolePostInfo{RoleInfo: role}
	this.ServeJSON()
}

func (this *RoleController) GetOne() {
	if _, e := this.ParseToken(); e != nil {
		this.RetError(e)
		return
	}

	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		this.RetError(errInputData)
		return
	}

	role := models.Role{}
	if code, err := role.FindById(id); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			this.RetError(errNoUser)
		} else {
			this.RetError(errDatabase)
		}
		return
	}
	beego.Debug("RoleInfo:", &role)

	role.ClearPass()

	this.Data["json"] = &models.RoleGetOneInfo{RoleInfo: &role}
	this.ServeJSON()
}

func (this *RoleController) GetAll() {
	if _, e := this.ParseToken(); e != nil {
		this.RetError(e)
		return
	}

	queryVal, queryOp, err := this.ParseQueryParm()
	if err != nil {
		beego.Debug("ParseQuery:", err)
		this.RetError(errInputData)
		return
	}
	beego.Debug("QueryVal:", queryVal)
	beego.Debug("QueryOp:", queryOp)

	order, err := this.ParseOrderParm()
	if err != nil {
		beego.Debug("ParseOrder:", err)
		this.RetError(errInputData)
		return
	}
	beego.Debug("Order:", order)

	limit, err := this.ParseLimitParm()
	/*
		if err != nil {
			beego.Debug("ParseLimit:", err)
			this.RetError(errInputData)
			return
		}
	*/
	beego.Debug("Limit:", limit)

	offset, err := this.ParseOffsetParm()
	/*
		if err != nil {
			beego.Debug("ParseOffset:", err)
			this.RetError(errInputData)
			return
		}
	*/
	beego.Debug("Offset:", offset)

	roles, err := models.GetAllRoles(queryVal, queryOp, order,
		limit, offset)
	if err != nil {
		beego.Error("GetAllRole:", err)
		this.RetError(errDatabase)
		return
	}
	beego.Debug("GetAllRole:", &roles)

	for i, _ := range roles {
		roles[i].ClearPass()
	}

	this.Data["json"] = &models.RoleGetAllInfo{RolesInfo: roles}
	this.ServeJSON()
}

func (this *RoleController) Put() {
	token, e := this.ParseToken()
	if e != nil {
		this.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		this.RetError(errPermission)
		return
	}

	idStr := this.Ctx.Input.Param(":id")
	if claims["id"] != idStr && claims["id"] != "1" {
		this.RetError(errPermission)
		return
	}

	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		this.RetError(errInputData)
		return
	}

	form := models.RolePutForm{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePut:", err)
		this.RetError(errInputData)
		return
	}
	beego.Debug("ParseRolePut:", &form)

	role := models.Role{}
	if code, err := role.UpdateById(id, &form); err != nil {
		beego.Error("UpdateRoleById:", err)
		this.RetError(errDatabase)
		return
	} else if code == models.ErrNotFound {
		this.RetError(errNoUserChange)
		return
	}

	if code, err := role.FindById(id); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			this.RetError(errNoUser)
		} else {
			this.RetError(errDatabase)
		}
		return
	}
	beego.Debug("NewRoleInfo:", &role)

	role.ClearPass()

	this.Data["json"] = &models.RolePutInfo{RoleInfo: &role}
	this.ServeJSON()
}

func (this *RoleController) Delete() {
	token, e := this.ParseToken()
	if e != nil {
		this.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		if claims["id"] != "1" {
			this.RetError(errPermission)
			return
		}
	} else {
		this.RetError(errPermission)
		return
	}

	idStr := this.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		this.RetError(errInputData)
		return
	}

	role := models.Role{}
	if code, err := role.DeleteById(id); err != nil {
		beego.Error("DeleteRoleById:", err)
		this.RetError(errDatabase)
		return
	} else if code == models.ErrNotFound {
		this.RetError(errNoUser)
		return
	}
}
