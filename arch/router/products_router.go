package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func ProductsRouter(ctx *echo.Echo, productsController controller.ProductController) {

	ctx.GET("/products", productsController.GetProducts)
	ctx.GET("products/:id", productsController.GetProductByID)
	ctx.POST("/product", productsController.CreateProducts)
	ctx.PATCH("product/:id", productsController.UpdateProducts)
	ctx.DELETE("product/:id", productsController.DeleteProduct)
}
