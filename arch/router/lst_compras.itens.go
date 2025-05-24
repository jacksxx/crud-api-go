package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func LstComprasItensRouter(ctx *echo.Echo, baseRouter string, lstComprasItensController controller.LstComprasItensController) {

	prodRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	prodRouter.GET("", lstComprasItensController.GetLstComprasItens)
	// prodRouter.GET("/:id", lstComprasItensController.GetProductByID)
	// prodRouter.POST("", lstComprasItensController.CreateProducts)
	// prodRouter.PATCH("/:id", lstComprasItensController.UpdateProducts)
	// prodRouter.PATCH("/inativar/:id", lstComprasItensController.InactivateProduct)
	// prodRouter.PATCH("/ativar/:id", lstComprasItensController.ActivateProduct)
}
