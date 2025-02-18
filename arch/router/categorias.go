package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func CategoriesRouter(ctx *echo.Echo, baseRouter string, categoriesController controller.CategoriaController) {

	catRouter := ctx.Group(baseRouter)

	catRouter.GET("", categoriesController.GetCategorias)
	catRouter.GET("/:id", categoriesController.GetCategoriaByID)
	catRouter.POST("", categoriesController.CreateCategoria)
	catRouter.PATCH("/:id", categoriesController.UpdateCategorias)
	catRouter.PATCH("/inativar/:id", categoriesController.InactivateCategorias)
	catRouter.PATCH("/ativar/:id", categoriesController.ActivateCategorias)

}
