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
