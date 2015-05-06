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
	var role models.Role
	err := json.Unmarshal(this.Ctx.Input.RequestBody, &role)
	if err != nil {
		beego.Debug("Input err: ", err)
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	role.RegDate = time.Now()
	beego.Debug("Input role: ", role)

	if code, err := role.Insert(); err != nil {
		beego.Debug("Insert role:", err)
		if code == 100 {
			this.Data["json"] = models.NewErrorInfo(ErrDupUser)
		} else {
			this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		}
		this.ServeJson()
		return
	}
	beego.Debug("Current role: ", role)

	this.Data["json"] = &role
	this.ServeJson()
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
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}
			var key string
			var value string
			var operator string
			if !nameRule.MatchString(kov[0]) {
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}
			key = kov[0]
			if op, ok := sqlOp[kov[1]]; ok {
				operator = op
			} else {
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}
			value = strings.Replace(kov[2], "'", "\\'", -1)

			queryVal[key] = value
			queryOp[key] = operator
		}
	}
	beego.Debug("queryVal", queryVal)
	beego.Debug("queryOp", queryOp)

	var order map[string]string = make(map[string]string)
	if v := this.GetString("order"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.Split(cond, ":")
			if len(kv) != 2 {
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}
			if !nameRule.MatchString(kv[0]) {
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}
			if kv[1] != "asc" && kv[1] != "desc" {
				this.Data["json"] =
					models.NewErrorInfo(ErrInputData)
				this.ServeJson()
				return
			}

			order[kv[0]] = kv[1]
		}
	}
	beego.Debug("order ", order)

	var limit int64 = 10
	if v, err := this.GetInt64("limit"); err != nil {
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	} else {
		if v > 0 {
			limit = v
		}
	}
	beego.Debug("limit ", limit)

	var offset int64 = 0
	if v, err := this.GetInt64("offset"); err != nil {
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	} else {
		if v > 0 {
			offset = v
		}
	}
	beego.Debug("offset ", offset)

	records, err := models.GetAllRoles(queryVal, queryOp, order,
		limit, offset)
	if err != nil {
		beego.Debug("get err ", err)
		this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		this.ServeJson()
		return
	}

	this.Data["json"] = &records
	this.ServeJson()
}

func (this *RoleController) Delete() {
	idStr := this.Ctx.Input.Params[":id"]
	id, err := strconv.ParseInt(idStr, 0, 64)
	if err != nil {
		this.Data["json"] = models.NewErrorInfo(ErrInputData)
		this.ServeJson()
		return
	}

	role := models.Role{}
	code, err := role.DeleteById(id)
	if err != nil {
		beego.Debug("DeleteRoleById:", err)
		this.Data["json"] = models.NewErrorInfo(ErrDatabase)
		this.ServeJson()
		return
	}

	if code == 100 {
		this.Data["json"] = models.NewErrorInfo(ErrNoUser)
	} else {
		this.Data["json"] = models.NewNormalInfo("Succes")
	}
	this.ServeJson()
}
