package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"crud-api-go/helper"
	"net/http"
	//"strconv"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"
)

type SubcategoriaController interface {
	GetSubcategorias(ctx echo.Context) error
	//GetSubcategoriaByID(ctx echo.Context) error
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

// func (c *subcategoriaController) GetSubcategoriaByID(ctx echo.Context) error {
// 	id := ctx.Param("id")
// 	if id == "" {
// 		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id é obrigatório"})
// 	}

// 	categoriaId, err := strconv.Atoi(id)
// 	if err != nil {
// 		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da categoria inválido"})
// 	}

// 	categoria, httpStatus, err := c.service.GetProductByID(categoriaId)
// 	if err != nil {
// 		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
// 	}

// 	return helper.BuildResponse(ctx, httpStatus, categoria, nil)
// }
