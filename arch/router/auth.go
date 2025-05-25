package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func AuthRoutes(e *echo.Echo, baseRouter string, controller controller.AuthController) {
	authrouter := e.Group(baseRouter)

	authrouter.POST("/login", controller.Login)
	authrouter.POST("/logout", controller.Logout)
	authrouter.POST("/refresh", controller.Refresh)
}
