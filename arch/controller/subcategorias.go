package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"crud-api-go/helper"
	"net/http"
	"strconv"

	//"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type SubcategoriaController interface {
	GetSubcategorias(ctx echo.Context) error
	GetSubcategoriaByID(ctx echo.Context) error
	CreateSubcategoria(ctx echo.Context) error
	UpdateSubcategorias(ctx echo.Context) error
	InactivateSubcategorias(ctx echo.Context) error
	ActivateSubcategorias(ctx echo.Context) error
}

type subcategoriaController struct {
	service    service.SubcategoriasService
	validate   *validator.Validate
	translator ut.Translator
}

func NewSubcategoriaController(service service.SubcategoriasService, validate *validator.Validate, translator ut.Translator) SubcategoriaController {
	return &subcategoriaController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

// GetSubcategorias godoc
// @Summary      Lista Subcategorias
// @Description  Retorna uma lista de subcategorias com suporte a filtros por nome e status.
// @Tags         Subcategorias
// @Accept       json
// @Produce      json
// @Param        nome    query     string false "Filtrar pelo nome (ILIKE)"
// @Param        status  query     string false "Filtrar pelo status"
// @Param        categorias_id query int  false "Filtrar por ID da categoria"
// @Param        limit   query     int    true  "Limite de itens por página (mínimo: 1)"
// @Param        page    query     int    true  "Número da página (mínimo: 1)"
// @Success      200     {object}  model.WebResponse{data=[]model.Subcategorias}
// @Failure      400     {object}  model.WebResponse
// @Failure      500     {object}  model.WebResponse
// @Router       /subcategorias [get]
func (c *subcategoriaController) GetSubcategorias(ctx echo.Context) error {
	filters := model.SubcategoriasFilters{}

	err := ctx.Bind(&filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	subcategorias, _, err := c.service.GetSubcategorias(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, subcategorias)
}

// GetSubcategoriaByID godoc
// @Summary Busca uma subcategoria pelo ID
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path int true "ID da Subcategoria"
// @Success 200 {object} model.WebResponse{data=model.Subcategorias}
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /subcategorias/{id} [get]
func (c *subcategoriaController) GetSubcategoriaByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id é obrigatório"})
	}

	subcategoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da categoria inválido"})
	}

	subcategoria, httpStatus, err := c.service.GetSubcategoriasById(subcategoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, subcategoria, nil)
}

// CreateSubcategoria godoc
// @Summary Cria uma nova subcategoria
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param subcategoria body model.SubcategoriasPost true "Dados da Subcategoria"
// @Success 201 {object} model.WebResponse{data=model.Subcategorias}
// @Failure 400 {object} model.WebResponse
// @Router /subcategorias [post]
func (c *subcategoriaController) CreateSubcategoria(ctx echo.Context) error {

	subcategoria := model.SubcategoriasPost{}

	validationErrors := helper.BindAndValidate(ctx, c.validate, &subcategoria, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	insertedSubcategories, httpStatus, err := c.service.CreateSubcategorias(subcategoria)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, insertedSubcategories, nil)
}

// UpdateSubcategorias godoc
// @Summary Atualiza uma subcategoria existente
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path int true "ID da Subcategoria"
// @Param subcategoria body model.SubcategoriasUpdate true "Dados da Subcategoria para atualização"
// @Success 200 {object} model.WebResponse{data=model.Subcategorias}
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /subcategorias/{id} [patch]
func (c *subcategoriaController) UpdateSubcategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	subcategoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	subcategoria := model.SubcategoriasUpdate{}
	validationErrors := helper.BindAndValidate(ctx, c.validate, &subcategoria, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	subcategoria.Id = subcategoriaId

	updatedSubcategoty, httpStatus, err := c.service.UpdateSubcategorias(subcategoria)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedSubcategoty, nil)
}

// InactivateSubcategorias godoc
// @Summary Inativa uma subcategoria
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path int true "ID da Subcategoria"
// @Success 204 {object} model.WebResponse
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /subcategorias/{id}/inativar [patch]
func (c *subcategoriaController) InactivateSubcategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	subcategoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = c.service.InactivateSubategorias(subcategoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}

// ActivateSubcategorias godoc
// @Summary Ativa uma subcategoria
// @Tags Subcategorias
// @Accept json
// @Produce json
// @Param id path int true "ID da Subcategoria"
// @Success 204 {object} model.WebResponse
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /subcategorias/{id}/ativar [patch]
func (c *subcategoriaController) ActivateSubcategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	subcategoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = c.service.ActivateSubategorias(subcategoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}
	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}
