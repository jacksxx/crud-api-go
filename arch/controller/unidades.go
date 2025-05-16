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

type UnitController interface {
	GetUnits(ctx echo.Context) error
	GetUnitByID(ctx echo.Context) error
	CreateUnit(ctx echo.Context) error
	UpdateUnit(ctx echo.Context) error
}

type unitController struct {
	service    service.UnidadesService
	validate   *validator.Validate
	translator ut.Translator
}

func NewUnitController(service service.UnidadesService, validate *validator.Validate, translator ut.Translator) UnitController {
	return &unitController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

// GetUnits retorna uma lista paginada de unidades com filtro por descrição
//
// @Summary      Lista unidades
// @Description  Retorna uma lista paginada de unidades com filtro por descrição
// @Tags         Unidades
// @Accept       json
// @Produce      json
// @Param        descricao  query     string  false  "Filtrar por descrição (ILIKE)"
// @Param        limit      query     int     true   "Quantidade de registros por página"
// @Param        page       query     int     true   "Número da página"
// @Success      200        {array}   model.Unidades
// @Failure      400        {object}  model.WebResponse
// @Failure      500        {object}  model.WebResponse
// @Router       /unidades [get]
func (c *unitController) GetUnits(ctx echo.Context) error {
	filters := model.UnidadesFilters{}

	err := ctx.Bind(&filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	categorias, _, err := c.service.GetUnits(filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	return helper.SuccessResponse(ctx, categorias)
}

// GetUnitByID retorna uma unidade específica pelo ID
//
// @Summary      Buscar unidade por ID
// @Description  Retorna uma unidade específica pelo ID
// @Tags         Unidades
// @Accept       json
// @Produce      json
// @Param        id   path      int     true   "ID da Unidade"
// @Success      200  {object}  model.Unidades
// @Failure      400  {object}  model.WebResponse
// @Failure      404  {object}  model.WebResponse
// @Router       /unidades/{id} [get]
func (c *unitController) GetUnitByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id é obrigatório"})
	}

	unitId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da unidade é inválido"})
	}

	unit, httpStatus, err := c.service.GetUnitByID(unitId)
	if err != nil {
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, unit, nil)
}

// CreateUnit cria uma nova unidade
//
// @Summary      Criar unidade
// @Description  Cria uma nova unidade
// @Tags         Unidades
// @Accept       json
// @Produce      json
// @Param        unidade  body      model.UnidadesPost  true  "Nova unidade"
// @Success      201      {object}  model.UnidadesPost
// @Failure      400      {object}  model.WebResponse
// @Router       /unidades [post]
func (c *unitController) CreateUnit(ctx echo.Context) error {

	unit := model.UnidadesPost{}

	validationErrors := helper.BindAndValidate(ctx, c.validate, &unit, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	insertedUnit, httpStatus, err := c.service.CreateUnit(unit)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	return helper.BuildResponse(ctx, httpStatus, insertedUnit, nil)
}

// UpdateUnit atualiza os dados de uma unidade existente
//
// @Summary      Atualizar unidade
// @Description  Atualiza os dados de uma unidade existente
// @Tags         Unidades
// @Accept       json
// @Produce      json
// @Param        id       path      int                  true  "ID da unidade"
// @Param        unidade  body      model.UnidadesUpdate true  "Dados da unidade a serem atualizados"
// @Success      200      {object}  model.UnidadesUpdate
// @Failure      400      {object}  model.WebResponse
// @Failure      500      {object}  model.WebResponse
// @Router       /unidades/{id} [patch]

func (c *unitController) UpdateUnit(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	unitId, err := strconv.Atoi(id)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	unit := model.UnidadesUpdate{}
	validationErrors := helper.BindAndValidate(ctx, c.validate, &unit, c.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	unit.Id = unitId

	updatedUnit, httpStatus, err := c.service.UpdateUnit(unit)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedUnit, nil)
}
