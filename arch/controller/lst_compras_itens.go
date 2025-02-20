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

type LstComprasItensController interface {
	GetLstComprasItens(ctx echo.Context) error
}

type lstComprasItensController struct {
	service    service.LstComprasItensService
	validate   *validator.Validate
	translator ut.Translator
}

func NewLstComprasItensController(service service.LstComprasItensService, validate *validator.Validate, translator ut.Translator) LstComprasItensController {
	return &lstComprasItensController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

func (c *lstComprasItensController) GetLstComprasItens(ctx echo.Context) error {
	filters := model.LstCompras_Itens_Filters{}

	err := ctx.Bind(&filters)

	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	Itens, err := c.service.GetLstComprasItens(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, Itens)
}
