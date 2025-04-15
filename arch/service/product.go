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
	InactivateProduct(id int) error
	ActivateProduct(id int) error
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
	tx, err := ps.repository.BeginTransaction()
	if err != nil {
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Chama o serviço para verificar se a categoria existe
	err = ps.repository.ValidateCategory(product.Categoria_Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar produto: %v", err)
	}
	// Chama o serviço para verificar se a subcategoria existe
	err = ps.repository.ValidateSubCategory(product.Categoria_Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da subcategoria não exista
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar produto: %v", err)
	}
	// Chama o repositório para criar um novo produto e retorna o ID gerado e o nome da categoria.
	insertedProduct, err := ps.repository.CreateProducts(product, tx)
	if err != nil {
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar produto: %v", err) // Retorna erro caso a inserção falhe.
	}

	err = tx.Commit()
	if err != nil {
		return model.ProductPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return insertedProduct, http.StatusOK, nil // Retorna o produto criado com ID e nome da categoria.
}

// UpdateProducts recebe um produto e atualiza seus dados no banco de dados.
func (ps *productService) UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, int, error) {
	tx, err := ps.repository.BeginTransaction()
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = ps.repository.ValidateProduct(product.Id, tx)
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err)
	}

	err = ps.repository.ValidateCategory(product.Categoria_Id, tx)
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err)
	}
	
	// Chama o serviço para verificar se a subcategoria existe
	err = ps.repository.ValidateSubCategory(product.Categoria_Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da subcategoria não exista
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err)
	}

	sts, err := ps.repository.VerificarStatusProduto(product.Id, tx)
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err)
	}
	if sts == "inativo" {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("o produto se encontra inativo")
	}

	updatedProduct, err := ps.repository.UpdateProducts(product, tx) // Chama o repositório para atualizar o produto.
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar produto: %v", err) // Retorna erro caso a atualização falhe.
	}

	err = tx.Commit()
	if err != nil {
		return model.ProductUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return updatedProduct, http.StatusOK, nil // Retorna o produto atualizado.
}

func (ps *productService) InactivateProduct(id int) error {
	tx, err := ps.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = ps.repository.ValidateProduct(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao inativar produto: %v", err)
	}

	sts, err := ps.repository.VerificarStatusProduto(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao inativar produto: %v", err)
	}
	if sts == "inativo" {
		return fmt.Errorf("o produto já se encontra inativo")
	}

	err = ps.repository.InactivateProduct(id, tx)
	if err != nil {
		return err // Retorna erro caso a exclusão falhe
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}

func (ps *productService) ActivateProduct(id int) error {
	tx, err := ps.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = ps.repository.ValidateProduct(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao ativar produto: %v", err)
	}

	sts, err := ps.repository.VerificarStatusProduto(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao inativar produto: %v", err)
	}
	if sts == "ativo" {
		return fmt.Errorf("o produto já se encontra ativo")
	}

	err = ps.repository.ActivateProduct(id, tx)
	if err != nil {
		return err // Retorna erro caso a exclusão falhe
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}
