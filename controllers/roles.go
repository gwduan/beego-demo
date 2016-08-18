package controllers

import (
	"encoding/json"
	"strconv"
	"time"

	"beego-demo/models"

	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
)

// RoleController definiton.
type RoleController struct {
	BaseController
}

// Auth method.
func (c *RoleController) Auth() {
	form := models.RoleAuthForm{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRoleAuth:", err)
		c.RetError(errInputData)
		return
	}
	beego.Debug("ParseRoleAuth:", &form)

	role := models.Role{}
	if code, err := role.FindByID(form.ID); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			c.RetError(errNoUser)
		} else {
			c.RetError(errDatabase)
		}
		return
	}
	beego.Debug("RoleInfo:", &role)

	if role.Name != form.Name || role.Password != form.Password {
		c.RetError(errPass)
		return
	}

	// Create the token with some claims
	claims := make(jwt.MapClaims)
	claims["id"] = strconv.FormatInt(form.ID, 10)
	claims["name"] = form.Name
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		beego.Error("jwt.SignedString:", err)
		c.RetError(errSystem)
		return
	}

	c.Data["json"] = &models.RoleAuthInfo{Token: tokenString}
	c.ServeJSON()
}

// Post method.
func (c *RoleController) Post() {
	token, e := c.ParseToken()
	if e != nil {
		c.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		if claims["id"] != "1" {
			c.RetError(errPermission)
			return
		}
	} else {
		c.RetError(errPermission)
		return
	}

	form := models.RolePostForm{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePost:", err)
		c.RetError(errInputData)
		return
	}
	beego.Debug("ParseRolePost:", &form)

	regDate := time.Now()
	role := models.NewRole(&form, regDate)
	beego.Debug("NewRole:", role)

	if code, err := role.Insert(); err != nil {
		beego.Error("InsertRole:", err)
		if code == models.ErrDupRows {
			c.RetError(errDupUser)
		} else {
			c.RetError(errDatabase)
		}
		return
	}

	role.ClearPass()

	c.Data["json"] = &models.RolePostInfo{RoleInfo: role}
	c.ServeJSON()
}

// GetOne method.
func (c *RoleController) GetOne() {
	if _, e := c.ParseToken(); e != nil {
		c.RetError(e)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		c.RetError(errInputData)
		return
	}

	role := models.Role{}
	if code, err := role.FindByID(id); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			c.RetError(errNoUser)
		} else {
			c.RetError(errDatabase)
		}
		return
	}
	beego.Debug("RoleInfo:", &role)

	role.ClearPass()

	c.Data["json"] = &models.RoleGetOneInfo{RoleInfo: &role}
	c.ServeJSON()
}

// GetAll method.
func (c *RoleController) GetAll() {
	if _, e := c.ParseToken(); e != nil {
		c.RetError(e)
		return
	}

	queryVal, queryOp, err := c.ParseQueryParm()
	if err != nil {
		beego.Debug("ParseQuery:", err)
		c.RetError(errInputData)
		return
	}
	beego.Debug("QueryVal:", queryVal)
	beego.Debug("QueryOp:", queryOp)

	order, err := c.ParseOrderParm()
	if err != nil {
		beego.Debug("ParseOrder:", err)
		c.RetError(errInputData)
		return
	}
	beego.Debug("Order:", order)

	limit, err := c.ParseLimitParm()
	/*
		if err != nil {
			beego.Debug("ParseLimit:", err)
			c.RetError(errInputData)
			return
		}
	*/
	beego.Debug("Limit:", limit)

	offset, err := c.ParseOffsetParm()
	/*
		if err != nil {
			beego.Debug("ParseOffset:", err)
			c.RetError(errInputData)
			return
		}
	*/
	beego.Debug("Offset:", offset)

	roles, err := models.GetAllRoles(queryVal, queryOp, order, limit, offset)
	if err != nil {
		beego.Error("GetAllRole:", err)
		c.RetError(errDatabase)
		return
	}
	beego.Debug("GetAllRole:", &roles)

	for i := range roles {
		roles[i].ClearPass()
	}

	c.Data["json"] = &models.RoleGetAllInfo{RolesInfo: roles}
	c.ServeJSON()
}

// Put method.
func (c *RoleController) Put() {
	token, e := c.ParseToken()
	if e != nil {
		c.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		c.RetError(errPermission)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	if claims["id"] != idStr && claims["id"] != "1" {
		c.RetError(errPermission)
		return
	}

	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		c.RetError(errInputData)
		return
	}

	form := models.RolePutForm{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePut:", err)
		c.RetError(errInputData)
		return
	}
	beego.Debug("ParseRolePut:", &form)

	role := models.Role{}
	if code, err := role.UpdateByID(id, &form); err != nil {
		beego.Error("UpdateRoleById:", err)
		c.RetError(errDatabase)
		return
	} else if code == models.ErrNotFound {
		c.RetError(errNoUserChange)
		return
	}

	if code, err := role.FindByID(id); err != nil {
		beego.Error("FindRoleById:", err)
		if code == models.ErrNotFound {
			c.RetError(errNoUser)
		} else {
			c.RetError(errDatabase)
		}
		return
	}
	beego.Debug("NewRoleInfo:", &role)

	role.ClearPass()

	c.Data["json"] = &models.RolePutInfo{RoleInfo: &role}
	c.ServeJSON()
}

// Delete method.
func (c *RoleController) Delete() {
	token, e := c.ParseToken()
	if e != nil {
		c.RetError(e)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok {
		if claims["id"] != "1" {
			c.RetError(errPermission)
			return
		}
	} else {
		c.RetError(errPermission)
		return
	}

	idStr := c.Ctx.Input.Param(":id")
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		c.RetError(errInputData)
		return
	}

	role := models.Role{}
	if code, err := role.DeleteByID(id); err != nil {
		beego.Error("DeleteRoleById:", err)
		c.RetError(errDatabase)
		return
	} else if code == models.ErrNotFound {
		c.RetError(errNoUser)
		return
	}
}
