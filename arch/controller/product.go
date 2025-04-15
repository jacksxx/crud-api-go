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

type ProductController interface {
	GetProducts(ctx echo.Context) error
	GetProductByID(ctx echo.Context) error
	CreateProducts(ctx echo.Context) error
	UpdateProducts(ctx echo.Context) error
	InactivateProduct(ctx echo.Context) error
	ActivateProduct(ctx echo.Context) error
}

// ProductController é responsável por lidar com as requisições HTTP relacionadas a produtos.
type productController struct {
	// ProductService contém a lógica de negócio e comunicação com o repositório.
	service    service.ProductService
	validate   *validator.Validate
	translator ut.Translator
}

// NewProductController cria e retorna uma nova instância de ProductController.
func NewProductController(service service.ProductService, validate *validator.Validate, translator ut.Translator) ProductController {
	return &productController{
		service:    service,
		validate:   validate,
		translator: translator,
	}
}

// GetProducts lida com requisições GET para buscar todos os produtos.

func (p *productController) GetProducts(ctx echo.Context) error {
	// Inicializa os filtros
	filters := model.ProductFilters{}

	// Faz o binding dos parâmetros da query para a struct `filters`
	err := ctx.Bind(&filters)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	// Chama o serviço para obter os produtos, passando os filtros
	products, _, err := p.service.GetProducts(filters)
	if err != nil {
		// Retorna erro 500 caso a busca falhe
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	// Retorna os produtos encontrados com status 200 (OK)
	return helper.SuccessResponse(ctx, products)
}

// GetProductByID lida com requisições GET para buscar um produto pelo ID.

func (pc *productController) GetProductByID(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id é obrigatório"})
	}
	// Converte o ID de string para inteiro.
	productId, err := strconv.Atoi(id)
	if err != nil {
		// Retorna erro 400 se o ID não for um número válido.
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"O Id da categoria inválido"})
	}
	// Chama o serviço para buscar o produto pelo ID.
	product, httpStatus, err := pc.service.GetProductByID(productId)
	if err != nil {
		// Retorna erro 500 caso ocorra um problema na busca.
		helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}
	// Retorna o produto encontrado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, product, nil)
}

// CreateProducts lida com requisições POST para criar um novo produto.
func (pc *productController) CreateProducts(ctx echo.Context) error {
	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	product := model.ProductPost{}

	validationErrors := helper.BindAndValidate(ctx, pc.validate, &product, pc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	// Chama o serviço para criar o produto no banco de dados.
	insertedProduct, httpStatus, err := pc.service.CreateProducts(product)
	if err != nil {
		// Retorna erro 500 caso ocorra uma falha na criação do produto.
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{err.Error()})
	}

	// Retorna o produto criado com status 201 (Created).
	return helper.BuildResponse(ctx, httpStatus, insertedProduct, nil)
}

// UpdateProducts lida com requisições PUT para atualizar um produto existente.
func (pc *productController) UpdateProducts(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro.
	productId, err := strconv.Atoi(id)
	if err != nil {
		// Retorna erro 400 se o ID não for um número válido.
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	product := model.ProductUpdate{}

	// Faz o bind dos dados da requisição para o objeto product.
	validationErrors := helper.BindAndValidate(ctx, pc.validate, &product, pc.translator)
	if len(validationErrors) > 0 {
		return helper.BuildResponse(ctx, http.StatusBadRequest, nil, validationErrors)
	}

	// Atribui o ID ao produto para garantir que estamos atualizando o item correto.
	product.Id = productId

	// Chama o serviço para atualizar o produto no banco de dados.
	updatedProduct, httpStatus, err := pc.service.UpdateProducts(product)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return helper.BuildResponse(ctx, httpStatus, nil, []string{err.Error()})
	}

	// Retorna o produto atualizado com status 200 (OK).
	return helper.BuildResponse(ctx, httpStatus, updatedProduct, nil)
}

func (pc *productController) InactivateProduct(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL
	id := ctx.Param("id")
	if id == "" {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = pc.service.InactivateProduct(productId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}

func (pc *productController) ActivateProduct(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL
	id := ctx.Param("id")
	if id == "" {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id não pode ser nulo"})
	}

	// Converte o ID de string para inteiro
	productId, err := strconv.Atoi(id)
	if err != nil {
		helper.BuildResponse(ctx, http.StatusBadRequest, nil, []string{"Id precisa ser um número"})
	}

	// Chama o serviço para excluir o produto
	err = pc.service.ActivateProduct(productId)
	if err != nil {
		return helper.BuildResponse(ctx, http.StatusInternalServerError, nil, []string{err.Error()})
	}

	// Retorna resposta de sucesso
	return helper.BuildResponse(ctx, http.StatusNoContent, nil, nil)
}
