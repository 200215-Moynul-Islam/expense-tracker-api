package controllers

import beego "github.com/beego/beego/v2/server/web"

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Success(status int, message string, data interface{}) {
	response := map[string]interface{}{
		"success": true,
		"message": message,
	}

	if data != nil {
		response["data"] = data
	}

	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *BaseController) Error(status int, message string) {
	c.Ctx.Output.SetStatus(status)

	c.Data["json"] = map[string]interface{}{
		"success": false,
		"message": message,
	}

	c.ServeJSON()
}