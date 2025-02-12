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
	DeleteCategorias(ctx echo.Context) error
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

func (cc *categoriaController) CreateCategoria(ctx echo.Context) error {

	categoria := model.CategoriasPost{}

	validationErrors := helper.BindAndValidate(ctx, cc.validate, &categoria, cc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	err := cc.service.ValidateCategoryName(categoria.Name, nil)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	insertedCategories, httpStatus, err := cc.service.CreateCategorias(categoria)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, insertedCategories, nil)
}

func (cc *categoriaController) UpdateCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para verificar se a categoria existe
	err = cc.service.ValidateCategory(categoriaId)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	categoria := model.CategoriasUpdate{}
	validationErrors := helper.BindAndValidate(ctx, cc.validate, &categoria, cc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	categoria.Id = categoriaId

	err = cc.service.ValidateCategoryName(categoria.Name, &categoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	updatedCategoty, httpStatus, err := cc.service.UpdateCategorias(categoria)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedCategoty, nil)
}

func (cc *categoriaController) DeleteCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para verificar se a categoria existe
	err = cc.service.ValidateCategory(categoriaId)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	// Chama o serviço para excluir o produto
	err = cc.service.DeleteCategorias(categoriaId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{"Erro ao deletar categoria"})
	}
	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}
