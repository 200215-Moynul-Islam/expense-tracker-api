package routers

import (
	"expense-tracker-api/controllers"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	ns := beego.NewNamespace("/api/v1",

		beego.NSRouter("/health", &controllers.HealthController{}),

		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
		),
	)

	beego.AddNamespace(ns)
}