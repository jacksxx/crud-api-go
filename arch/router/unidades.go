package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func UnitsRouter(ctx *echo.Echo, baseRouter string, unitController controller.UnitController) {

	unitRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	unitRouter.GET("", unitController.GetUnits)
	unitRouter.GET("/:id", unitController.GetUnitByID)
	unitRouter.POST("", unitController.CreateUnit)
	unitRouter.PATCH("/:id", unitController.UpdateUnit)

}
