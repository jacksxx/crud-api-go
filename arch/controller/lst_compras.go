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

type LstComprasController interface {
	GetLstCompras(ctx echo.Context) error
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
