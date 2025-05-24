package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func ProductsRouter(ctx *echo.Echo, baseRouter string, productsController controller.ProductController) {

	prodRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	prodRouter.GET("", productsController.GetProducts)
	prodRouter.GET("/:id", productsController.GetProductByID)
	prodRouter.POST("", productsController.CreateProducts)
	prodRouter.PATCH("/:id", productsController.UpdateProducts)
	prodRouter.PATCH("/inativar/:id", productsController.InactivateProduct)
	prodRouter.PATCH("/ativar/:id", productsController.ActivateProduct)
}
