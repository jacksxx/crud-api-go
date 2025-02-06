package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductController struct {
	ProductService service.ProductService
}

func NewProductController(service service.ProductService) ProductController {
	return ProductController{
		ProductService: service,
	}
}

func (p *ProductController) GetProducts(ctx echo.Context) error {
	products, err := p.ProductService.GetProducts()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)

}

func (pc *ProductController) GetProductByID(ctx echo.Context) error {
	id := ctx.Param("id")
	if id == "" {
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	productId, err := strconv.Atoi(id)
	if err != nil {
		rsp := model.ResponseMessage("Id precisa ser um número")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	product, err := pc.ProductService.GetProductByID(productId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if product == nil {
		rsp := model.ResponseMessage("Produto não encontrado")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}
	
	return ctx.JSON(http.StatusOK, product)
}

func (pc *ProductController) CreateProducts(ctx echo.Context) error {
	var product model.Product

	err := ctx.Bind(&product)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, err)
	}
	insertedProduct, err := pc.ProductService.CreateProducts(product)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusCreated, insertedProduct)

}
