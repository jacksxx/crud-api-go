package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
)

type ProductService interface {
	GetProducts(filters model.ProductFilters) ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	CreateProducts(product model.ProductPost) (model.ProductPost, error)
	UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error)
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
func (p *productService) GetProducts(filters model.ProductFilters) ([]model.Product, error) {
	return p.repository.GetProducts(filters)
}

// GetProductByID busca um produto pelo ID no repositório.
func (ps *productService) GetProductByID(id int) (*model.Product, error) {
	product, err := ps.repository.GetProductByID(id)
	if err != nil {
		return nil, err // Retorna erro caso a consulta falhe.
	}
	return product, nil // Retorna o produto encontrado.
}

// CreateProducts recebe um produto e o insere no banco de dados.
func (ps *productService) CreateProducts(product model.ProductPost) (model.ProductPost, error) {
	// Chama o repositório para criar um novo produto e retorna o ID gerado e o nome da categoria.
	productId, categoriaName, err := ps.repository.CreateProducts(product)
	if err != nil {
		return model.ProductPost{}, err // Retorna erro caso a inserção falhe.
	}

	// Atualiza os valores do produto criado.
	product.Id = productId
	product.Categoria_Name = categoriaName // Atualiza o nome da categoria

	return product, nil // Retorna o produto criado com ID e nome da categoria.
}

// UpdateProducts recebe um produto e atualiza seus dados no banco de dados.
func (ps *productService) UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error) {
	updatedProduct, err := ps.repository.UpdateProducts(product) // Chama o repositório para atualizar o produto.
	if err != nil {
		return model.ProductUpdate{}, err // Retorna erro caso a atualização falhe.
	}

	return updatedProduct, nil // Retorna o produto atualizado.
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
