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
	beego.Controller
}

func (this *RoleController) Post() {
	form := models.RolePostForm{}
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePost:", err)
		this.CustomAbort(errInputData.Ret())
		return
	}
	beego.Debug("ParseRolePost:", &form)

	regDate := time.Now()
	role := models.NewRole(&form, regDate)
	beego.Debug("NewRole:", role)

	if code, err := role.Insert(); err != nil {
		beego.Debug("InsertRole:", err)
		if code == 100 {
			this.CustomAbort(errDupUser.Ret())
		} else {
			this.CustomAbort(errDatabase.Ret())
		}
		return
	}

	//this.Data["json"] = models.NewNormalInfo("Succes")
	//this.ServeJson()
}

func (this *RoleController) GetOne() {
	idStr := this.Ctx.Input.Params[":id"]
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		beego.Debug("ParseRoleId:", err)
		this.CustomAbort(errInputData.Ret())
		return
	}

	role := models.Role{}
	if code, err := role.FindById(id); err != nil {
		beego.Debug("FindRoleById:", err)
		if code == 100 {
			this.CustomAbort(errNoUser.Ret())
		} else {
			this.CustomAbort(errDatabase.Ret())
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
				this.CustomAbort(errInputData.Ret())
				return
			}
			var key string
			var value string
			var operator string
			if !nameRule.MatchString(kov[0]) {
				this.CustomAbort(errInputData.Ret())
				return
			}
			key = kov[0]
			if op, ok := sqlOp[kov[1]]; ok {
				operator = op
			} else {
				this.CustomAbort(errInputData.Ret())
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
				this.CustomAbort(errInputData.Ret())
				return
			}
			if !nameRule.MatchString(kv[0]) {
				this.CustomAbort(errInputData.Ret())
				return
			}
			if kv[1] != "asc" && kv[1] != "desc" {
				this.CustomAbort(errInputData.Ret())
				return
			}

			order[kv[0]] = kv[1]
		}
	}
	beego.Debug("Order:", order)

	var limit int64 = 10
	if v, err := this.GetInt64("limit"); err != nil {
		this.CustomAbort(errInputData.Ret())
		return
	} else {
		if v > 0 {
			limit = v
		}
	}
	beego.Debug("Limit:", limit)

	var offset int64 = 0
	if v, err := this.GetInt64("offset"); err != nil {
		this.CustomAbort(errInputData.Ret())
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
		beego.Debug("GetAllRole:", err)
		this.CustomAbort(errDatabase.Ret())
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
		this.CustomAbort(errInputData.Ret())
		return
	}

	form := models.RolePutForm{}
	err = json.Unmarshal(this.Ctx.Input.RequestBody, &form)
	if err != nil {
		beego.Debug("ParseRolePut:", err)
		this.CustomAbort(errInputData.Ret())
		return
	}
	beego.Debug("ParseRolePut:", &form)

	role := models.Role{}
	if code, err := role.UpdateById(id, &form); err != nil {
		beego.Debug("UpdateRoleById:", err)
		this.CustomAbort(errDatabase.Ret())
		return
	} else if code == 100 {
		this.CustomAbort(errNoUserChange.Ret())
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
		this.CustomAbort(errInputData.Ret())
		return
	}

	role := models.Role{}
	if code, err := role.DeleteById(id); err != nil {
		beego.Debug("DeleteRoleById:", err)
		this.CustomAbort(errDatabase.Ret())
		return
	} else if code == 100 {
		this.CustomAbort(errNoUser.Ret())
		return
	}

	//this.Data["json"] = models.NewNormalInfo("Succes")
	//this.ServeJson()
}
