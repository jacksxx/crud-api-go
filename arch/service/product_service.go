package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
)

// ProductService é responsável por interagir com o repositório e aplicar regras de negócio.
type ProductService struct {
	// ProductRepository é a camada que interage diretamente com o banco de dados.
	ProductRepository repository.ProductRepository
}

// NewProductService inicializa um novo serviço de produtos com um repositório injetado.
func NewProductService(repository repository.ProductRepository) ProductService {
	return ProductService{
		ProductRepository: repository,
	}
}

// GetProducts busca todos os produtos do repositório e retorna uma lista de produtos.
func (p *ProductService) GetProducts(filters model.ProductFilters) ([]model.Product, error) {
	return p.ProductRepository.GetProducts(filters)
}

// GetProductByID busca um produto pelo ID no repositório.
func (ps *ProductService) GetProductByID(id int) (*model.Product, error) {
	product, err := ps.ProductRepository.GetProductByID(id)
	if err != nil {
		return nil, err // Retorna erro caso a consulta falhe.
	}
	return product, nil // Retorna o produto encontrado.
}

// CreateProducts recebe um produto e o insere no banco de dados.
func (ps *ProductService) CreateProducts(product model.Product) (model.Product, error) {
	// Chama o repositório para criar um novo produto e retorna o ID gerado e o nome da categoria.
	productId, categoriaName, err := ps.ProductRepository.CreateProducts(product)
	if err != nil {
		return model.Product{}, err // Retorna erro caso a inserção falhe.
	}

	// Atualiza os valores do produto criado.
	product.Id = productId
	product.Categoria_Name = categoriaName // Atualiza o nome da categoria

	return product, nil // Retorna o produto criado com ID e nome da categoria.
}

// UpdateProducts recebe um produto e atualiza seus dados no banco de dados.
func (ps *ProductService) UpdateProducts(product model.Product) (model.Product, error) {
	updatedProduct, err := ps.ProductRepository.UpdateProducts(product) // Chama o repositório para atualizar o produto.
	if err != nil {
		return model.Product{}, err // Retorna erro caso a atualização falhe.
	}

	return updatedProduct, nil // Retorna o produto atualizado.
}

func (ps *ProductService) DeleteProduct(id int) error {
	// Chama o repositório para excluir o produto
	err := ps.ProductRepository.DeleteProduct(id)
	if err != nil {
		return err // Retorna erro caso a exclusão falhe
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}
