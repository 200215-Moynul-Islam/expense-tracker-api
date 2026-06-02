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
			beego.NSInclude(
				&controllers.AuthController{},
			),
			beego.NSRouter("/register", &controllers.AuthController{}, "post:Register"),
			beego.NSRouter("/login", &controllers.AuthController{}, "post:Login"),
		),

		beego.NSNamespace("/expenses",
			beego.NSBefore(middlewares.RequireAuthentication),

			beego.NSInclude(
				&controllers.ExpenseController{},
			),

			beego.NSRouter("/summary", &controllers.ExpenseController{}, "get:GetSummary"),

			beego.NSRouter("", &controllers.ExpenseController{}, "get:GetAll"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "get:GetByID"),
			beego.NSRouter("", &controllers.ExpenseController{}, "post:Create"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "put:Update"),
			beego.NSRouter("/:id", &controllers.ExpenseController{}, "delete:Delete"),
		),
	)

	beego.AddNamespace(ns)
}
