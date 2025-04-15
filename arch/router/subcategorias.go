package router

import (
	"crud-api-go/arch/controller"

	"github.com/labstack/echo/v4"
)

func SubcategoriesRouter(ctx *echo.Echo, baseRouter string, subcategoriesController controller.SubcategoriaController) {

	catRouter := ctx.Group(baseRouter)

	catRouter.GET("", subcategoriesController.GetSubcategorias)
	// catRouter.GET("/:id", categoriesController.GetCategoriaByID)
	// catRouter.POST("", categoriesController.CreateCategoria)
	// catRouter.PATCH("/:id", categoriesController.UpdateCategorias)
	// catRouter.PATCH("/inativar/:id", categoriesController.InactivateCategorias)
	// catRouter.PATCH("/ativar/:id", categoriesController.ActivateCategorias)

}
