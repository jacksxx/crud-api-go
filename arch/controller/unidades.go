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
