package controllers

import "net/http"

type HealthController struct {
	BaseController
}

// Get checks server health
// @Title Health Check
// @Description Check whether the server is running
// @Tags Health
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{} "Server is running"
// @router /health [get]
func (c *HealthController) Get() {
	c.Success(http.StatusOK, "Server is running", nil)
}