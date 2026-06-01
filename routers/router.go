package routers

import (
	"expense-tracker-api/controllers"
	"expense-tracker-api/middlewares"

	beego "github.com/beego/beego/v2/server/web"
)

func init() {

	ns := beego.NewNamespace("/api/v1",

		beego.NSRouter("/health", &controllers.HealthController{}),

		beego.NSNamespace("/auth",
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),

		beego.NSNamespace("/expenses",
			beego.NSBefore(middlewares.RequireAuthentication),

			beego.NSRouter("", &controllers.ExpenseController{}, "post:Create"),
		),
	)

	beego.AddNamespace(ns)
}