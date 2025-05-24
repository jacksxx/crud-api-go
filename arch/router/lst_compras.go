package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func LstComprasRouter(ctx *echo.Echo, baseRouter string, lstComprasController controller.LstComprasController) {

	prodRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	prodRouter.GET("", lstComprasController.GetLstCompras)
	prodRouter.GET("/:id", lstComprasController.GetLstComprasByCodigo)
	prodRouter.POST("", lstComprasController.PostLstCompras)
	prodRouter.PATCH("/:id", lstComprasController.UpdateLstCompras)
	prodRouter.PATCH("/cancel/:id", lstComprasController.CancelLstCompras)
	prodRouter.PATCH("/:id/finalizar", lstComprasController.FinishLstCompras)
}
