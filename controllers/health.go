package controllers

type HealthController struct {
	BaseController
}

func (c *HealthController) Get() {
	c.Success(200, "Server is running", nil)
}