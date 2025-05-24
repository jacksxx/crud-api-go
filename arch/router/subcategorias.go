package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func SubcategoriesRouter(ctx *echo.Echo, baseRouter string, subcategoriesController controller.SubcategoriaController) {

	catRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	catRouter.GET("", subcategoriesController.GetSubcategorias)
	catRouter.GET("/:id", subcategoriesController.GetSubcategoriaByID)
	catRouter.POST("", subcategoriesController.CreateSubcategoria)
	catRouter.PATCH("/:id", subcategoriesController.UpdateSubcategorias)
	catRouter.PATCH("/inativar/:id", subcategoriesController.InactivateSubcategorias)
	catRouter.PATCH("/ativar/:id", subcategoriesController.ActivateSubcategorias)

}
