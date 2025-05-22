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

// GetLstComprasItens
// @Summary Lista os itens de uma lista de compras
// @Description Retorna os itens agrupados por lista e por status (Comprados / Não Comprados)
// @Tags Itens da Lista de Compras
// @Accept json
// @Produce json
// @Param lst_compras_id query int true "ID da lista de compras"
// @Param products_id query int false "ID do produto"
// @Param products_name query string false "Nome do produto (busca parcial)"
// @Param page query int false "Página" minimum(1)
// @Param limit query int false "Limite de itens por página" minimum(1)
// @Success 200 {object} model.WebResponse{data=map[int]map[string][]model.LstCompras_Itens}
// @Failure 400 {object} model.WebResponse
// @Failure 500 {object} model.WebResponse
// @Router /lst_compras_itens [get]
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
