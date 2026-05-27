package controllers

import beego "github.com/beego/beego/v2/server/web"

type HealthController struct {
	beego.Controller
}

func (c *HealthController) Get() {
	c.Data["json"] = map[string]interface{}{
		"success": true,
		"message": "Server is running",
	}
	c.ServeJSON()
}