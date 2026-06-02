package controllers

import beego "github.com/beego/beego/v2/server/web"

type APIResponse struct {
	Success bool `json:"success"`
	Message string `json:"message"`
	Data interface{} `json:"data,omitempty"`
}

type BaseController struct {
	beego.Controller
}

func (c *BaseController) Success(status int, message string, data interface{}) {
	response := APIResponse{
		Success: true,
		Message: message,
		Data: data,
	}

	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = response
	c.ServeJSON()
}

func (c *BaseController) Error(status int, message string) {
	response := APIResponse{
		Success: false,
		Message: message,
	}

	c.Ctx.Output.SetStatus(status)
	c.Data["json"] = response
	c.ServeJSON()
}