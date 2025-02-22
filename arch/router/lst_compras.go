package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func LstComprasRouter(ctx *echo.Echo, baseRouter string, lstComprasController controller.LstComprasController) {

	prodRouter := ctx.Group(baseRouter)

	prodRouter.GET("", lstComprasController.GetLstCompras)
	prodRouter.GET("/:id", lstComprasController.GetLstComprasByCodigo)
	prodRouter.POST("", lstComprasController.PostAluguel)
	// prodRouter.PATCH("/:id", lstComprasItensController.UpdateProducts)
	// prodRouter.PATCH("/inativar/:id", lstComprasItensController.InactivateProduct)
	// prodRouter.PATCH("/ativar/:id", lstComprasItensController.ActivateProduct)
}
