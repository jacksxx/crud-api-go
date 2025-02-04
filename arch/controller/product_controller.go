package controller

import (
	"crud-api-go/arch/service"
	"net/http"

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
		ctx.JSON(http.StatusInternalServerError, err)
	}
	return ctx.JSON(http.StatusOK, products)

}
