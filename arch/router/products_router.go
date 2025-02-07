package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func ProductsRouter(ctx *echo.Echo, baseRouter string, productsController controller.ProductController) {

	prodRouter := ctx.Group(baseRouter)

	prodRouter.GET("", productsController.GetProducts)
	prodRouter.GET("/:id", productsController.GetProductByID)
	prodRouter.POST("", productsController.CreateProducts)
	prodRouter.PATCH("/:id", productsController.UpdateProducts)
	prodRouter.DELETE("/:id", productsController.DeleteProduct)
}
