package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

// ProductRepository representa a estrutura do repositório de produtos, armazenando a conexão com o banco de dados
type ProductRepository struct {
	connection *sql.DB
}

// NewProductRepository cria uma nova instância do ProductRepository com a conexão do banco de dados
func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

// GetProducts busca todos os produtos no banco de dados e retorna uma lista de produtos
func (pr *ProductRepository) GetProducts(filters model.ProductFilters) ([]model.Product, error) {
	// Inicia a query base
	query := `SELECT products_id, products_name, products_price FROM prod.products`
	var conditions []string
	var args []interface{}
	argIndex := 1

	// Adiciona filtro por nome, se fornecido
	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("products_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	// Adiciona as condições à query, se houver filtros
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Adiciona a ordenação por nome ascendente
	query += " ORDER BY products_name ASC"

	// Aplicar paginação (LIMIT e OFFSET)
	if filters.Limit > 0 {
		query += fmt.Sprintf(" LIMIT $%d", argIndex)
		args = append(args, filters.Limit)
		argIndex++
	}
	if filters.Page > 0 {
		offset := (filters.Page - 1) * filters.Limit
		query += fmt.Sprintf(" OFFSET $%d", argIndex)
		args = append(args, offset)
	}

	// Executa a consulta no banco de dados
	rows, err := pr.connection.Query(query, args...)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	defer rows.Close() // Garante que o cursor do banco será fechado após o uso

	// Lista para armazenar os produtos retornados
	var productList []model.Product

	// Itera sobre os resultados da consulta
	for rows.Next() {
		var product model.Product
		err = rows.Scan(&product.Id, &product.Name, &product.Price)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
		productList = append(productList, product)
	}

	return productList, nil
}

// GetProductByID busca um produto específico pelo ID no banco de dados
func (pr *ProductRepository) GetProductByID(id int) (*model.Product, error) {
	// Prepara a query SQL para evitar SQL Injection
	query, err := pr.connection.Prepare("SELECT * FROM prod.products WHERE products_id = $1")
	if err != nil {
		fmt.Println(err) // Log do erro
		return nil, err  // Retorna erro caso a preparação da query falhe
	}

	var product model.Product
	// Executa a query e faz o scan do resultado para o objeto product
	err = query.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		// Verifica se o erro é do tipo "nenhuma linha encontrada"
		if err == sql.ErrNoRows {
			return nil, nil // Retorna nil para indicar que o produto não foi encontrado
		}
		return nil, err // Retorna erro caso ocorra outro tipo de falha
	}
	query.Close()        // Fecha a consulta preparada
	return &product, nil // Retorna o produto encontrado
}

// CreateProducts insere um novo produto no banco de dados e retorna o ID do produto inserido
func (pr *ProductRepository) CreateProducts(product model.Product) (int, error) {
	var Id int
	// Prepara a query de inserção para evitar SQL Injection
	query, err := pr.connection.Prepare("INSERT INTO prod.products" + "(products_name, products_price)" + "VALUES ($1,$2) RETURNING products_id")
	if err != nil {
		fmt.Println(err) // Log do erro
		return 0, err    // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	// Executa a query e escaneia o ID do novo produto inserido
	err = query.QueryRow(product.Name, product.Price).Scan(&Id)
	if err != nil {
		fmt.Println(err) // Log do erro
		return 0, err    // Retorna erro caso ocorra uma falha na inserção
	}
	return Id, nil // Retorna o ID do novo produto inserido
}

// UpdateProducts atualiza um produto existente no banco de dados
func (pr *ProductRepository) UpdateProducts(product model.Product) (model.Product, error) {
	// Prepara a query de atualização para evitar SQL Injection
	query, err := pr.connection.Prepare("UPDATE prod.products SET products_name = $1, products_price = $2 WHERE products_id = $3 RETURNING products_id")
	if err != nil {
		fmt.Println(err)            // Log do erro
		return model.Product{}, err // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	var Id int
	// Executa a query e escaneia o ID do produto atualizado
	err = query.QueryRow(product.Name, product.Price, product.Id).Scan(&Id)
	if err != nil {
		fmt.Println(err)            // Log do erro
		return model.Product{}, err // Retorna erro caso ocorra uma falha na atualização
	}
	return product, nil // Retorna o produto atualizado
}

func (pr *ProductRepository) DeleteProduct(id int) error {
	// Prepara a query de exclusão para evitar SQL Injection
	query, err := pr.connection.Prepare("DELETE FROM prod.products WHERE products_id = $1")
	if err != nil {
		fmt.Println(err) // Log do erro
		return err       // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	// Executa a query de exclusão
	_, err = query.Exec(id)
	if err != nil {
		fmt.Println(err) // Log do erro
		return err       // Retorna erro caso ocorra uma falha na exclusão
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}
