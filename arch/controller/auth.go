package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"crud-api-go/helper"
	"net/http"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

// AuthController é o controlador para autenticação e gerenciamento de tokens
type AuthController struct {
	service service.AuthService
	validate    *validator.Validate
	translator  ut.Translator
}

// NewAuthController cria um novo controlador de autenticação
func NewAuthController(service service.AuthService, validate *validator.Validate, translator ut.Translator) AuthController {
	return AuthController{
		service: service,
		validate:    validate,
		translator:  translator,
	}
}


// Login
// @Summary Autentica um usuário
// @Description Autentica um usuário com base nas credenciais fornecidas
// @Tags Autenticação
// @Accept json
// @Produce json
// @Param credentials body model.LoginCredentials true "Credenciais do usuário"
// @Success 200 {object} model.WebResponse "Tokens de autenticação"
// @Failure 400 {object} model.WebResponse "Erro de validação"
// @Failure 400 {object} model.WebResponse "Campos de usuário e senha são obrigatórios"
// @Failure 500 {object} model.WebResponse "usuário não encontrado"
// @Failure 401 {object} model.WebResponse "Credenciais inválidas"
// @Failure 500 {object} model.WebResponse "Erro interno do servidor"
// @Router /auth/login [post]
func (u *AuthController) Login(ctx echo.Context) error {
	var credentials model.LoginCredentials

	if err := ctx.Bind(&credentials); err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}
	if credentials.Username == "" || credentials.Password == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Campos de usuário e senha são obrigatórios"})
	}

	httpStatus, err := u.service.Authenticate(ctx, credentials)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, http.StatusOK, []string{"Autenticado com sucesso"}, nil)
}

// Refresh
// @Summary Atualiza o token de acesso
// @Description Atualiza o token de acesso do usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Success 200 {object} model.WebResponse "Token atualizado com sucesso"
// @Failure 400 {object} model.WebResponse "Token de atualização não encontrado"
// @Failure 403 {object} model.WebResponse "Sessão expirada. Por favor, faça login novamente"
// @Failure 500 {object} model.WebResponse "Erro interno do servidor"
// @Router /auth/refresh [post]
func (u *AuthController) Refresh(ctx echo.Context) error {
	_, httpStatus, err := u.service.Refresh(ctx)

	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, http.StatusOK, []string{"Token atualizado com sucesso"}, nil)
}

// Logout é o endpoint para deslogar um usuário
// @Summary Desloga um usuário
// @Description Desloga um usuário
// @Tags Autenticação
// @Accept json
// @Produce json
// @Success 200 {object} model.WebResponse "Logout bem-sucedido"
// @Failure 500 {object} model.WebResponse "Erro interno do servidor"
// @Router /auth/logout [post]
func (u *AuthController) Logout(ctx echo.Context) error {
	err := u.service.Logout(ctx)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	return helper.BuildResponse(ctx, http.StatusOK, []string{"Logout bem-sucedido"}, nil)
}
