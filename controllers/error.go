package controllers

// ErrorController definiton.
type ErrorController struct {
	BaseController
}

// RetError return error informatino in JSON.
func (c *ErrorController) RetError(e *ControllerError) {
	c.Data["json"] = e
	c.ServeJSON()
}

// Error404 redefine 404 error information.
func (c *ErrorController) Error404() {
	c.RetError(err404)
}
