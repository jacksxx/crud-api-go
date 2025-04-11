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
	PostLstCompras(ctx echo.Context) error
	UpdateLstCompras(ctx echo.Context) error
	CancelLstCompras(ctx echo.Context) error
	FinishLstCompras(ctx echo.Context) error
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

func (c *lstComprasController) PostLstCompras(ctx echo.Context) error {
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

func (c *lstComprasController) UpdateLstCompras(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	lst_compras := model.LstCompras_Update{}
	// Faz o bind dos dados da requisição para o objeto product.
	validationErrors := helper.BindAndValidate(ctx, c.validate, &lst_compras, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}
	// Converte o ID de string para inteiro.
	lst_compras_Id, err := strconv.Atoi(id)
	if err != nil || lst_compras_Id <= 0 {
		// Retorna erro 400 se o ID não for um número válido.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}
	// Atribui o ID ao produto para garantir que estamos atualizando o item correto.
	lst_compras.Id = lst_compras_Id

	// Chama o serviço para atualizar o produto no banco de dados.
	updatedProduct, httpStatus, err := c.service.UpdateLstCompras(lst_compras)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedProduct, nil)
}

func (c *lstComprasController) FinishLstCompras(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	lst_compras := model.LstCompras_Finish{}
	// Faz o bind dos dados da requisição para o objeto product.
	validationErrors := helper.BindAndValidate(ctx, c.validate, &lst_compras, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}
	// Converte o ID de string para inteiro.
	lst_compras_Id, err := strconv.Atoi(id)
	if err != nil || lst_compras_Id <= 0 {
		// Retorna erro 400 se o ID não for um número válido.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}
	// Atribui o ID ao produto para garantir que estamos atualizando o item correto.
	lst_compras.Id = lst_compras_Id

	// Chama o serviço para atualizar o produto no banco de dados.
	finishedLst, httpStatus, err := c.service.FinishLstCompras(lst_compras)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, finishedLst, nil)
}

func (c *lstComprasController) CancelLstCompras(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	lst_compras := model.LstCompras_Cancel{}
	validationErrors := helper.BindAndValidate(ctx, c.validate, &lst_compras, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	lst_compras_Id, err := strconv.Atoi(id)
	if err != nil || lst_compras_Id <= 0 {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}
	lst_compras.Id = lst_compras_Id

	httpStatus, err := c.service.CancelLstCompras(lst_compras)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, nil, nil)
}
