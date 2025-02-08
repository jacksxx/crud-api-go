package controller

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/service"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ProductController interface {
	GetProducts(ctx echo.Context) error
	GetProductByID(ctx echo.Context) error
	CreateProducts(ctx echo.Context) error
	UpdateProducts(ctx echo.Context) error
	DeleteProduct(ctx echo.Context) error
}

// ProductController é responsável por lidar com as requisições HTTP relacionadas a produtos.
type productController struct {
	// ProductService contém a lógica de negócio e comunicação com o repositório.
	service service.ProductService
}

// NewProductController cria e retorna uma nova instância de ProductController.
func NewProductController(service service.ProductService) ProductController {
	return &productController{
		service: service,
	}
}

// GetProducts lida com requisições GET para buscar todos os produtos.

func (p *productController) GetProducts(ctx echo.Context) error {
	// Inicializa os filtros
	var filters model.ProductFilters

	// Faz o binding dos parâmetros da query para a struct `filters`
	err := ctx.Bind(&filters)
	if err != nil {
		return ctx.JSON(http.StatusBadRequest, map[string]string{"error": "Parâmetros inválidos"})
	}

	// Chama o serviço para obter os produtos, passando os filtros
	products, err := p.service.GetProducts(filters)
	if err != nil {
		// Retorna erro 500 caso a busca falhe
		return ctx.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Retorna os produtos encontrados com status 200 (OK)
	return ctx.JSON(http.StatusOK, products)
}

// GetProductByID lida com requisições GET para buscar um produto pelo ID.

func (pc *productController) GetProductByID(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}
	// Converte o ID de string para inteiro.
	productId, err := strconv.Atoi(id)
	if err != nil {
		// Retorna erro 400 se o ID não for um número válido.
		rsp := model.ResponseMessage("Id precisa ser um número")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}
	// Chama o serviço para buscar o produto pelo ID.
	product, err := pc.service.GetProductByID(productId)
	if err != nil {
		// Retorna erro 500 caso ocorra um problema na busca.
		return ctx.JSON(http.StatusInternalServerError, err)
	}
	if product == nil {
		// Retorna erro 400 se o produto não for encontrado.
		rsp := model.ResponseMessage("Produto não encontrado")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}
	// Retorna o produto encontrado com status 200 (OK).
	return ctx.JSON(http.StatusOK, product)
}

// CreateProducts lida com requisições POST para criar um novo produto.
func (pc *productController) CreateProducts(ctx echo.Context) error {
	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	var product model.ProductPost

	// Faz o bind dos dados da requisição para o objeto product.
	err := ctx.Bind(&product)
	if err != nil {
		// Retorna erro 400 caso o bind falhe.
		return ctx.JSON(http.StatusBadRequest, err)
	}

	// Chama o serviço para verificar se a categoria existe
	err = pc.service.ValidateCategory(product.Categoria_Id)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage(err.Error()))
	}

	// Chama o serviço para criar o produto no banco de dados.
	insertedProduct, err := pc.service.CreateProducts(product)
	if err != nil {
		// Retorna erro 500 caso ocorra uma falha na criação do produto.
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	// Retorna o produto criado com status 201 (Created).
	return ctx.JSON(http.StatusCreated, insertedProduct)
}

// UpdateProducts lida com requisições PUT para atualizar um produto existente.
func (pc *productController) UpdateProducts(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL.
	id := ctx.Param("id")
	if id == "" {
		// Retorna erro 400 se o ID estiver vazio.
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Converte o ID de string para inteiro.
	productId, err := strconv.Atoi(id)
	if err != nil {
		// Retorna erro 400 se o ID não for um número válido.
		rsp := model.ResponseMessage("Id precisa ser um número")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Cria uma variável do tipo Product para armazenar os dados recebidos no corpo da requisição.
	var product model.ProductUpdate

	// Faz o bind dos dados da requisição para o objeto product.
	err = ctx.Bind(&product)
	if err != nil {
		// Retorna erro 400 caso o bind falhe.
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage("Erro ao processar os dados"))
	}

	// Atribui o ID ao produto para garantir que estamos atualizando o item correto.
	product.Id = productId

	// Chama o serviço para verificar se a categoria existe
	err = pc.service.ValidateCategory(product.Categoria_Id)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return ctx.JSON(http.StatusBadRequest, model.ResponseMessage(err.Error()))
	}

	// Chama o serviço para atualizar o produto no banco de dados.
	updatedProduct, err := pc.service.UpdateProducts(product)
	if err != nil {
		// Retorna erro 500 caso a atualização falhe.
		return ctx.JSON(http.StatusInternalServerError, model.ResponseMessage(err.Error()))
	}

	// Retorna o produto atualizado com status 200 (OK).
	return ctx.JSON(http.StatusOK, updatedProduct)
}

func (pc *productController) DeleteProduct(ctx echo.Context) error {
	// Obtém o ID do produto a partir dos parâmetros da URL
	id := ctx.Param("id")
	if id == "" {
		rsp := model.ResponseMessage("Id não pode ser nulo")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Converte o ID de string para inteiro
	productId, err := strconv.Atoi(id)
	if err != nil {
		rsp := model.ResponseMessage("Id precisa ser um número válido")
		return ctx.JSON(http.StatusBadRequest, rsp)
	}

	// Chama o serviço para excluir o produto
	err = pc.service.DeleteProduct(productId)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, model.ResponseMessage("Erro ao deletar produto"))
	}

	// Retorna resposta de sucesso
	return ctx.JSON(http.StatusNoContent, nil)
}
