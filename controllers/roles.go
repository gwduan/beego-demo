package controllers

import (
	"beego-demo/models"
	"encoding/json"
	"github.com/astaxie/beego"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RoleController struct {
	BaseController
}

func (this *RoleController) Post() {
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

	//this.Ctx.ResponseWriter.WriteHeader(201)
	this.Data["json"] = &models.RolePostInfo{RoleInfo: role}
	this.ServeJson()
}

func (this *RoleController) GetOne() {
	idStr := this.Ctx.Input.Params[":id"]
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
	this.ServeJson()
}

var sqlOp = map[string]string{
	"eq": "=",
	"ne": "<>",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
}

func (this *RoleController) GetAll() {
	var nameRule = regexp.MustCompile("^[a-zA-Z0-9_]+$")

	var queryVal map[string]string = make(map[string]string)
	var queryOp map[string]string = make(map[string]string)
	if v := this.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kov := strings.Split(cond, ":")
			if len(kov) != 3 {
				this.RetError(errInputData)
				return
			}
			var key string
			var value string
			var operator string
			if !nameRule.MatchString(kov[0]) {
				this.RetError(errInputData)
				return
			}
			key = kov[0]
			if op, ok := sqlOp[kov[1]]; ok {
				operator = op
			} else {
				this.RetError(errInputData)
				return
			}
			value = strings.Replace(kov[2], "'", "\\'", -1)

			queryVal[key] = value
			queryOp[key] = operator
		}
	}
	beego.Debug("QueryVal:", queryVal)
	beego.Debug("QueryOp:", queryOp)

	var order map[string]string = make(map[string]string)
	if v := this.GetString("order"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.Split(cond, ":")
			if len(kv) != 2 {
				this.RetError(errInputData)
				return
			}
			if !nameRule.MatchString(kv[0]) {
				this.RetError(errInputData)
				return
			}
			if kv[1] != "asc" && kv[1] != "desc" {
				this.RetError(errInputData)
				return
			}

			order[kv[0]] = kv[1]
		}
	}
	beego.Debug("Order:", order)

	var limit int64 = 10
	if v, err := this.GetInt64("limit"); err != nil {
		this.RetError(errInputData)
		return
	} else {
		if v > 0 {
			limit = v
		}
	}
	beego.Debug("Limit:", limit)

	var offset int64 = 0
	if v, err := this.GetInt64("offset"); err != nil {
		this.RetError(errInputData)
		return
	} else {
		if v > 0 {
			offset = v
		}
	}
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
	this.ServeJson()
}

func (this *RoleController) Put() {
	idStr := this.Ctx.Input.Params[":id"]
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

	//this.Data["json"] = models.NewNormalInfo("Succes")
	//this.ServeJson()
}

func (this *RoleController) Delete() {
	idStr := this.Ctx.Input.Params[":id"]
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

	//this.Data["json"] = models.NewNormalInfo("Succes")
	//this.ServeJson()
}
