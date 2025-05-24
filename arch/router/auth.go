package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, baseRouter string, controller controller.AuthController) {
	e.POST(baseRouter+"/login", controller.Login)
	e.POST(baseRouter+"/logout", controller.Logout)
	e.POST(baseRouter+"/refresh", controller.Refresh)
}