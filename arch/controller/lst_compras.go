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

type LstComprasController interface {
	GetLstCompras(ctx echo.Context) error
	GetLstComprasByCodigo(ctx echo.Context) error
	PostAluguel(ctx echo.Context) error
}

type lstComprasController struct {
	service    service.LstComprasService
	validate   *validator.Validate
	translator ut.Translator
}

func NewLstComprasController(service service.LstComprasService, validate *validator.Validate, translator ut.Translator) LstComprasController {
	return &lstComprasController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

func (c *lstComprasController) GetLstCompras(ctx echo.Context) error {
	filters := model.LstCompras_Filters{}

	err := ctx.Bind(&filters)

	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	Itens, _, err := c.service.GetLstCompras(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, Itens)
}

func (c *lstComprasController) GetLstComprasByCodigo(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Código é obrigatório"})
	}
	codigoInt, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Código do aluguel inválido"})
	}
	aluguel, httpStatus, err := c.service.GetLstComprasById(codigoInt)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}
	return helper.BuildResponse(ctx, httpStatus, aluguel, nil)
}

func (c *lstComprasController) PostAluguel(ctx echo.Context) error {
	compras := model.LstCompras_Post{}

	validationsErrors := helper.BindAndValidate(ctx, c.validate, &compras, c.translator)
	if len(validationsErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationsErrors)
	}
	comprasCriada, httpStatus, err := c.service.CreateLstCompras(compras)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, comprasCriada, nil)
}
