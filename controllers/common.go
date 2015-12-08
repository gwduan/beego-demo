package controllers

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/gwduan/beego"
	"github.com/gwduan/beego/validation"
	"regexp"
	"strings"
)

const (
	ErrInputData    = "数据输入错误"
	ErrDatabase     = "数据库操作错误"
	ErrDupUser      = "用户信息已存在"
	ErrNoUser       = "用户信息不存在"
	ErrPass         = "密码不正确"
	ErrNoUserPass   = "用户信息不存在或密码不正确"
	ErrNoUserChange = "用户信息不存在或数据未改变"
	ErrInvalidUser  = "用户信息不正确"
	ErrOpenFile     = "打开文件出错"
	ErrWriteFile    = "写文件出错"
	ErrSystem       = "操作系统错误"
)

type ControllerError struct {
	Status   int    `json:"status"`
	Code     int    `json:"code"`
	Message  string `json:"message"`
	DevInfo  string `json:"dev_info"`
	MoreInfo string `json:"more_info"`
}

var (
	errInputData    = &ControllerError{400, 10001, "数据输入错误", "客户端参数错误", ""}
	errDatabase     = &ControllerError{500, 10002, "服务器错误", "数据库操作错误", ""}
	errDupUser      = &ControllerError{400, 10003, "用户信息已存在", "数据库记录重复", ""}
	errNoUser       = &ControllerError{400, 10004, "用户信息不存在", "数据库记录不存在", ""}
	errPass         = &ControllerError{400, 10005, "用户信息不存在或密码不正确", "密码不正确", ""}
	errNoUserPass   = &ControllerError{400, 10006, "用户信息不存在或密码不正确", "数据库记录不存在或密码不正确", ""}
	errNoUserChange = &ControllerError{400, 10007, "用户信息不存在或数据未改变", "数据库记录不存在或数据未改变", ""}
	errInvalidUser  = &ControllerError{400, 10008, "用户信息不正确", "Session信息不正确", ""}
	errOpenFile     = &ControllerError{500, 10009, "服务器错误", "打开文件出错", ""}
	errWriteFile    = &ControllerError{500, 10010, "服务器错误", "写文件出错", ""}
	errSystem       = &ControllerError{500, 10011, "服务器错误", "操作系统错误", ""}
	errExpired      = &ControllerError{400, 10012, "登录已过期", "验证token过期", ""}
	errPermission   = &ControllerError{400, 10013, "没有权限", "没有操作权限", ""}
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) RetError(e *ControllerError) {
	if mode := beego.AppConfig.String("runmode"); mode == "prod" {
		e.DevInfo = ""
	}

	this.Ctx.Output.Header("Content-Type", "application/json; charset=utf-8")
	this.Ctx.ResponseWriter.WriteHeader(e.Status)
	this.Data["json"] = e
	this.ServeJSON()

	this.StopRun()
}

var sqlOp = map[string]string{
	"eq": "=",
	"ne": "<>",
	"gt": ">",
	"ge": ">=",
	"lt": "<",
	"le": "<=",
}

func (this *BaseController) ParseQueryParm() (v map[string]string, o map[string]string, err error) {
	var nameRule = regexp.MustCompile("^[a-zA-Z0-9_]+$")
	var queryVal map[string]string = make(map[string]string)
	var queryOp map[string]string = make(map[string]string)

	query := this.GetString("query")
	if query == "" {
		return queryVal, queryOp, nil
	}

	for _, cond := range strings.Split(query, ",") {
		kov := strings.Split(cond, ":")
		if len(kov) != 3 {
			return queryVal, queryOp,
				errors.New("Query format != k:o:v")
		}

		var key string
		var value string
		var operator string
		if !nameRule.MatchString(kov[0]) {
			return queryVal, queryOp,
				errors.New("Query key format is wrong")
		}
		key = kov[0]
		if op, ok := sqlOp[kov[1]]; ok {
			operator = op
		} else {
			return queryVal, queryOp,
				errors.New("Query operator is wrong")
		}
		value = strings.Replace(kov[2], "'", "\\'", -1)

		queryVal[key] = value
		queryOp[key] = operator
	}

	return queryVal, queryOp, nil
}

func (this *BaseController) ParseOrderParm() (o map[string]string, err error) {
	var nameRule = regexp.MustCompile("^[a-zA-Z0-9_]+$")
	var order map[string]string = make(map[string]string)

	v := this.GetString("order")
	if v == "" {
		return order, nil
	}

	for _, cond := range strings.Split(v, ",") {
		kv := strings.Split(cond, ":")
		if len(kv) != 2 {
			return order, errors.New("Order format != k:v")
		}
		if !nameRule.MatchString(kv[0]) {
			return order, errors.New("Order key format is wrong")
		}
		if kv[1] != "asc" && kv[1] != "desc" {
			return order, errors.New("Order val isn't asc/desc")
		}

		order[kv[0]] = kv[1]
	}

	return order, nil
}

func (this *BaseController) ParseLimitParm() (l int64, err error) {
	if v, err := this.GetInt64("limit"); err != nil {
		return 10, err
	} else if v > 0 {
		return v, nil
	} else {
		return 10, nil
	}
}

func (this *BaseController) ParseOffsetParm() (o int64, err error) {
	if v, err := this.GetInt64("offset"); err != nil {
		return 0, err
	} else if v > 0 {
		return v, nil
	} else {
		return 0, nil
	}
}

func (this *BaseController) VerifyForm(obj interface{}) (err error) {
	valid := validation.Validation{}
	ok, err := valid.Valid(obj)
	if err != nil {
		return err
	}
	if !ok {
		str := ""
		for _, err := range valid.Errors {
			str += err.Key + ":" + err.Message + ";"
		}
		return errors.New(str)
	}

	return nil
}

func (this *BaseController) ParseToken() (t *jwt.Token, e *ControllerError) {
	authString := this.Ctx.Input.Header("Authorization")
	beego.Debug("AuthString:", authString)

	kv := strings.Split(authString, " ")
	if len(kv) != 2 || kv[0] != "Bearer" {
		beego.Error("AuthString invalid:", authString)
		return nil, errInputData
	}
	tokenString := kv[1]

	// Parse token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil {
		beego.Error("Parse token:", err)
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// That's not even a token
				return nil, errInputData
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				// Token is either expired or not active yet
				return nil, errExpired
			} else {
				// Couldn't handle this token
				return nil, errInputData
			}
		} else {
			// Couldn't handle this token
			return nil, errInputData
		}
	}
	if !token.Valid {
		beego.Error("Token invalid:", tokenString)
		return nil, errInputData
	}
	beego.Debug("Token:", token)

	return token, nil
}
