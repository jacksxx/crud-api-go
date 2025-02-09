package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"net/http"
	"strconv"

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
	service service.CategoriaService
}

func NewCategoriaController(service service.CategoriaService) CategoriaController {
	return &categoriaController{
		service: service,
	}
}

func (cc *categoriaController) GetCategorias(ctx echo.Context) error {
	var filters model.CategoriasFilters

	err := ctx.Bind(&filters)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error:": "Parâmetros inválidos"})
	}

	categorias, err := cc.service.GetCategorias(filters)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error:": err.Error()})
	}

	return ctx.JSON(http.StatusOK, categorias)
}

func (cc *categoriaController) GetCategoriaByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		rsp := model.ResponseMessage("Id precisa ser um número")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	categoria, err := cc.service.GetProductByID(categoriaId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if categoria == nil {
		rsp := model.ResponseMessage("Categoria não encontrada")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	return ctx.JSON(http.StatusOK, categoria)
}

func (cc *categoriaController) CreateCategoria(ctx echo.Context) error {

	var categoria model.CategoriasPost

	err := ctx.Bind(&categoria)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}

	err = cc.service.ValidateCategoryName(categoria.Name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage(err.Error()))
	}

	insertedCategories, err := cc.service.CreateCategorias(categoria)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, insertedCategories)
}

func (cc *categoriaController) UpdateCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		rsp := model.ResponseMessage("Id precisa ser um número")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	var categoria model.CategoriasUpdate
	err = ctx.Bind(&categoria)
	if err != nil {
		// Retorna erro 400 caso o bind falhe.
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage("Erro ao processar os dados"))
	}

	categoria.Id = categoriaId

	err = cc.service.ValidateCategoryName(categoria.Name)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage(err.Error()))
	}

	updatedCategoty, err := cc.service.UpdateCategorias(categoria)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return ctx.JSON(http.StatusInternalServerError, model.ResponseMessage(err.Error()))
	}

	// Retorna o produto atualizado com status 200 (OK).
	return ctx.JSON(http.StatusOK, updatedCategoty)
}

func (cc *categoriaController) DeleteCategorias(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Converte o ID de string para inteiro
	categoriaId, err := strconv.Atoi(id)
	if err != nil {
		rsp := model.ResponseMessage("Id precisa ser um número válido")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Chama o serviço para excluir o produto
	err = cc.service.DeleteCategorias(categoriaId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ResponseMessage("Erro ao deletar categoria"))
	}

	// Retorna resposta de sucesso
	return ctx.JSON(http.StatusNoContent, nil)
}
