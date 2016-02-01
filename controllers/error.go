package controllers

type ErrorController struct {
	BaseController
}

func (this *ErrorController) RetError(e *ControllerError) {
	this.Data["json"] = e
	this.ServeJSON()
}

func (this *ErrorController) Error404() {
	this.RetError(err404)
}
