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

type CategoriaController interface {
	GetCategorias(ctx echo.Context) error
	GetCategoriaByID(ctx echo.Context) error
	CreateCategoria(ctx echo.Context) error
	UpdateCategorias(ctx echo.Context) error
	InactivateCategorias(ctx echo.Context) error
	ActivateCategorias(ctx echo.Context) error
}

type categoriaController struct {
	service    service.CategoriaService
	validate   *validator.Validate
	translator ut.Translator
}

func NewCategoriaController(service service.CategoriaService, validate *validator.Validate, translator ut.Translator) CategoriaController {
	return &categoriaController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

// GetCategorias godoc
// @Summary Lista categorias com filtros
// @Description Retorna uma lista de categorias com filtros opcionais por nome e status
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param nome query string false "Filtrar por nome (ILIKE)"
// @Param status query string false "Filtrar por status"
// @Param limit query int false "Limite de itens por página (mínimo: 1)"
// @Param page query int false "Número da página (mínimo: 1)"
// @Success 200 {object} model.WebResponse{data=[]model.Categorias}
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /categories [get]
func (cc *categoriaController) GetCategorias(ctx echo.Context) error {
	filters := model.CategoriasFilters{}

	err := ctx.Bind(&filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	categorias, _, err := cc.service.GetCategorias(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, categorias)
}

// GetCategoriaByID godoc
// @Summary Busca uma categoria por ID
// @Description Retorna os dados de uma categoria específica
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param id path int true "ID da categoria"
// @Success 200 {object} model.WebResponse{data=model.Categorias}
// @Failure 400 {object} model.WebResponse
// @Failure 404 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /categories/{id} [get]
func (cc *categoriaController) GetCategoriaByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id é obrigatório"})
	}

	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da categoria inválido"})
	}

	categoria, httpStatus, err := cc.service.GetProductByID(categoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, categoria, nil)
}

// CreateCategoria godoc
// @Summary Cria uma nova categoria
// @Description Cria uma nova categoria no sistema
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param categoria body model.CategoriasPost true "Dados da nova categoria"
// @Success 201 {object} model.WebResponse{data=model.Categorias}
// @Failure 400 {object} model.WebResponse
// @Router /categories [post]
func (cc *categoriaController) CreateCategoria(ctx echo.Context) error {

	categoria := model.CategoriasPost{}

	validationErrors := helper.BindAndValidate(ctx, cc.validate, &categoria, cc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	insertedCategories, httpStatus, err := cc.service.CreateCategorias(categoria)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, insertedCategories, nil)
}

// UpdateCategorias godoc
// @Summary Atualiza uma categoria existente
// @Description Atualiza os dados de uma categoria com base no ID
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param id path int true "ID da categoria"
// @Param categoria body model.CategoriasUpdate true "Dados atualizados da categoria"
// @Success 200 {object} model.WebResponse{data=model.Categorias}
// @Failure 400 {object} model.WebResponse
// @Failure 404 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /categories/{id} [patch]
func (cc *categoriaController) UpdateCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	categoria := model.CategoriasUpdate{}
	validationErrors := helper.BindAndValidate(ctx, cc.validate, &categoria, cc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	categoria.Id = categoriaId

	updatedCategoty, httpStatus, err := cc.service.UpdateCategorias(categoria)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedCategoty, nil)
}

// InactivateCategorias godoc
// @Summary Inativa uma categoria
// @Description Define data_inativacao e status como inativo para a categoria
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param id path int true "ID da categoria"
// @Success 204 {object} nil
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /categories/inativar/{id} [patch]
func (cc *categoriaController) InactivateCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = cc.service.InactivateCategorias(categoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}

// ActivateCategorias godoc
// @Summary Ativa uma categoria
// @Description Remove a data_inativacao e define status como ativo
// @Tags Categorias
// @Accept  json
// @Produce  json
// @Param id path int true "ID da categoria"
// @Success 204 {object} nil
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /categories/ativar/{id} [patch]
func (cc *categoriaController) ActivateCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = cc.service.ActivateCategorias(categoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}
