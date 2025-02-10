package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
	"net/http"
)

type ProductService interface {
	GetProducts(filters model.ProductFilters) (model.PaginatedResponse[model.Product], int, error)
	GetProductByID(id int) (model.Product, int, error)
	CreateProducts(product model.ProductPost) (model.ProductPost, int, error)
	UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, int, error)
	DeleteProduct(id int) error
	ValidateCategory(categoriaId int) error
}

// ProductService é responsável por interagir com o repositório e aplicar regras de negócio.
type productService struct {
	// ProductRepository é a camada que interage diretamente com o banco de dados.
	repository repository.ProductRepository
}

// NewProductService inicializa um novo serviço de produtos com um repositório injetado.
func NewProductService(repository repository.ProductRepository) ProductService {
	return &productService{
		repository: repository,
	}
}

// GetProducts busca todos os produtos do repositório e retorna uma lista de produtos.
func (ps *productService) GetProducts(filters model.ProductFilters) (model.PaginatedResponse[model.Product], int, error) {
	products, err := ps.repository.GetProducts(filters)
	if err != nil {
		return model.PaginatedResponse[model.Product]{}, 0, fmt.Errorf("erro ao buscar produtos: %v", err)
	}
	total, err := ps.repository.CountProducts(filters)
	if err != nil {
		return model.PaginatedResponse[model.Product]{}, 0, fmt.Errorf("erro ao buscar quantidade de produtos: %v", err)
	}
	response := model.PaginatedResponse[model.Product]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       products,
	}
	return response, total, nil
}

// GetProductByID busca um produto pelo ID no repositório.
func (ps *productService) GetProductByID(id int) (model.Product, int, error) {
	product, err := ps.repository.GetProductByID(id)
	if err != nil {
		return model.Product{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar produtos: %v", err) // Retorna erro caso a consulta falhe.
	}
	if product.Id == 0 {
		return model.Product{}, http.StatusInternalServerError, fmt.Errorf("produto não existe: %v", err)
	}
	return product, http.StatusOK, nil // Retorna o produto encontrado.
}

// CreateProducts recebe um produto e o insere no banco de dados.
func (ps *productService) CreateProducts(product model.ProductPost) (model.ProductPost, int, error) {
	// Chama o repositório para criar um novo produto e retorna o ID gerado e o nome da categoria.
	productId, categoriaName, err := ps.repository.CreateProducts(product)
	if err != nil {
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar produto: %v", err) // Retorna erro caso a inserção falhe.
	}

	// Atualiza os valores do produto criado.
	product.Id = productId
	product.Categoria_Name = categoriaName // Atualiza o nome da categoria

	return product, http.StatusOK, nil // Retorna o produto criado com ID e nome da categoria.
}

// UpdateProducts recebe um produto e atualiza seus dados no banco de dados.
func (ps *productService) UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, int, error) {
	updatedProduct, err := ps.repository.UpdateProducts(product) // Chama o repositório para atualizar o produto.
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err) // Retorna erro caso a atualização falhe.
	}

	return updatedProduct, http.StatusOK, nil // Retorna o produto atualizado.
}

func (ps *productService) DeleteProduct(id int) error {
	// Chama o repositório para excluir o produto
	err := ps.repository.DeleteProduct(id)
	if err != nil {
		return err // Retorna erro caso a exclusão falhe
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}

func (ps *productService) ValidateCategory(categoriaId int) error {
	err := ps.repository.ValidateCategory(categoriaId)
	if err != nil {
		return err
	}
	return nil
}
