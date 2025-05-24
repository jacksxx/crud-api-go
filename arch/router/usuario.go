package router

import (
	"crud-api-go/arch/controller"
	"crud-api-go/middleware"

	"github.com/labstack/echo/v4"
)

func UsuarioRouter(ctx *echo.Echo, baseRouter string, usersController controller.UsuarioController) {

	userRouter := ctx.Group(baseRouter, middleware.AuthMiddleware())

	userRouter.GET("", usersController.GetUsuario)
	userRouter.GET("/:id", usersController.GetUsuarioById)
	userRouter.POST("", usersController.CreateUsuario)
	userRouter.PATCH("/:id", usersController.UpdateUsuario)
	userRouter.PATCH("/inativar/:id", usersController.InactivateUsuario)
	userRouter.PATCH("/ativar/:id", usersController.ActivateUsuario)
	userRouter.PATCH("/resetpassword/:id", usersController.ResetarSenha)
	userRouter.PATCH("/password/:id", usersController.AtualizarSenha)

}
