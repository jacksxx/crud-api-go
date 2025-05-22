package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"crud-api-go/helper"
	"net/http"
	"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type UsuarioController interface {
	GetUsuario(ctx echo.Context) error
	GetUsuarioById(ctx echo.Context) error
	CreateUsuario(ctx echo.Context) error
	UpdateUsuario(ctx echo.Context) error
	InactivateUsuario(ctx echo.Context) error
	ActivateUsuario(ctx echo.Context) error
	AtualizarSenha(ctx echo.Context) error
	ResetarSenha(ctx echo.Context) error
}

type usuarioController struct {
	service    service.UsuarioService
	validate   *validator.Validate
	translator ut.Translator
}

func NewUsuarioController(service service.UsuarioService, validate *validator.Validate, translator ut.Translator) UsuarioController {
	return &usuarioController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

// GetUsuarios
// @Summary Lista usuarios com filtros
// @Description Retorna uma lista de usuarios com filtros opcionais
// @Tags Usuários
// @Accept  json
// @Produce  json
// @Param ativo query bool false "Filtrar por status"
// @Param limit query int false "Limite de itens por página (mínimo: 1)"
// @Param page query int false "Número da página (mínimo: 1)"
// @Success 200 {object} model.WebResponse{data=[]model.Usuario}
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /users [get]
func (c *usuarioController) GetUsuario(ctx echo.Context) error {
	filters := model.UsuarioFilter{}

	err := ctx.Bind(&filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	usuarios, err := c.service.GetUsuarios(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, usuarios)
}

// GetUsuarioById
// @Description Retorna os dados de um usuário específico pelo código
// @Summary Obtém um usuário por código
// @Tags Usuários
// @Produce json
// @Param id path int true "Código do usuário"
// @Success 200 {object} model.WebResponse{data=model.Usuario}
// @Failure 400 {object} model.WebResponse "Código é obrigatório"
// @Failure 404 {object} model.WebResponse "Usuário não encontrado"
// @Router /users/{id} [get]
func (c *usuarioController) GetUsuarioById(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Código é obrigatório"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da categoria inválido"})
	}

	usuario, httpStatus, err := c.service.GetUsuarioById(usuarioId)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}
	return helper.BuildResponse(ctx, httpStatus, usuario, nil)
}

// CreateUsuario
// @Description Cria um novo usuário com os dados fornecidos.
// @Summary Cria um novo usuário
// @Tags Usuários
// @Accept json
// @Produce json
// @Param usuario body model.UsuarioPost true "Dados do usuário"
// @Success 201 {object} model.WebResponse{data=model.Usuario}
// @Failure 400 {object} model.WebResponse "Erro na validação dos dados"
// @Failure 409 {object} model.WebResponse "Já existe um usuário com esse email ou nome de usuário"
// @Failure 500 {object} model.WebResponse "Erro interno do servidor"
// @Router /users [post]
func (c *usuarioController) CreateUsuario(ctx echo.Context) error {
	usuario := model.UsuarioPost{}

	if errors := helper.BindAndValidate(ctx, c.validate, &usuario, c.translator); errors != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, errors)
	}

	if validationErrors := usuario.Validate(); len(validationErrors) > 0 {
		errorMessages := make([]string, len(validationErrors))
		for i, err := range validationErrors {
			errorMessages[i] = err.Error()
		}
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, errorMessages)
	}

	usuarioCriado, httpStatus, err := c.service.CreateUsuario(usuario)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, usuarioCriado, nil)
}

// PatchUsuario atualiza um usuário
// @Description Atualiza os dados de um usuário específico
// @Summary Atualiza um usuário
// @Tags Usuários
// @Accept json
// @Produce json
// @Param id path int true "Código do usuário"
// @Param usuario body model.UsuarioUpdate true "Dados do usuário"
// @Success 200 {object} model.WebResponse{data=model.Usuario}
// @Failure 400 {object} model.WebResponse "Erro na validação dos dados"
// @Failure 404 {object} model.WebResponse "Usuário não encontrado"
// @Router /users/{id} [patch]
func (c *usuarioController) UpdateUsuario(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	usuario := model.UsuarioUpdate{}
	validationErrors := helper.BindAndValidate(ctx, c.validate, &usuario, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	if validationErrors := usuario.Validate(); len(validationErrors) > 0 {
		errorMessages := make([]string, len(validationErrors))
		for i, err := range validationErrors {
			errorMessages[i] = err.Error()
		}
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, errorMessages)
	}

	usuario.Id = usuarioId

	updatedCategoty, httpStatus, err := c.service.UpdateUsuario(usuario)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, updatedCategoty, nil)
}

// InativarUsuario
// @Description Inativa um usuário específico pelo código
// @Summary Inativa um usuário
// @Tags Usuários
// @Produce json
// @Param id path int true "Código do usuário"
// @Success 204 {object} model.WebResponse
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /users/inativar/{id} [patch]
func (c *usuarioController) InactivateUsuario(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	err = c.service.InativarUsuario(usuarioId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}

// AtivarUsuario
// @Description Ativa um usuário específico pelo código
// @Summary Ativa um usuário
// @Tags Usuários
// @Produce json
// @Param id path int true "Código do usuário"
// @Success 204 {object} model.WebResponse
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /users/ativar/{id} [patch]
func (c *usuarioController) ActivateUsuario(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	err = c.service.AtivarUsuario(usuarioId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}

// AtualizarSenha
// @Description Atualiza a senha de um usuário específico
// @Summary Atualiza a senha de um usuário
// @Tags usuários
// @Accept json
// @Produce json
// @Param id path int true "Código do usuário"
// @Param senha body model.SenhaUpdateRequest true "Nova senha"
// @Success 200 {object} model.WebResponse "Senha atualizada com sucesso"
// @Failure 400 {object} model.WebResponse "Erro na validação dos dados"
// @Failure 401 {object} model.WebResponse "Usuário não autorizado"
// @Failure 404 {object} model.WebResponse "Usuário não encontrado"
// @Router /users/password/{id} [patch]
func (c *usuarioController) AtualizarSenha(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	dados := model.SenhaUpdateRequest{}
	if errors := helper.BindAndValidate(ctx, c.validate, &dados, c.translator); errors != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, errors)
	}
	codigoUsuario, err := helper.ExtrairCodigoUsuarioDoToken(ctx)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusUnauthorized, nil, []string{err.Error()})
	}

	if codigoUsuario != usuarioId {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Você não tem permissão para atualizar a senha de outro usuário"})
	}

	httpStatus, err := c.service.UpdateSenha(usuarioId, dados)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, nil, nil)
}

// ResetarSenha reseta a senha de um usuário
// @Summary Reseta a senha de um usuário
// @Description Reseta a senha de um usuário para uma senha padrão (primeiro nome em minúsculas + ano atual)
// @Tags Usuários
// @Accept json
// @Produce json
// @Param codigo path int true "Código do usuário"
// @Success 200 {object} model.WebResponse
// @Failure 404 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /users/resetpassword/{id} [patch]
func (c *usuarioController) ResetarSenha(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	usuarioId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	httpStatus, err := c.service.ResetarSenha(usuarioId)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, nil, nil)
}
