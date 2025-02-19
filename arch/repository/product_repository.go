package repository

import (
	"crud-api-go/arch/model"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"strings"
)

// ProductRepository representa a estrutura do repositório de produtos, armazenando a conexão com o banco de dados
type ProductRepository interface {
	GetProducts(filters model.ProductFilters) ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)
	CreateProducts(product model.ProductPost) (int, string, error)
	UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error)
	InactivateProduct(id int) error
	ActivateProduct(id int) error
	ValidateCategory(categoriaId int) error
	CountProducts(filters model.ProductFilters) (int, error)
	ValidateProduct(productId int) error
}

// ProductRepository representa a estrutura do repositório de produtos, armazenando a conexão com o banco de dados
type productRepository struct {
	connection *sql.DB
}

// NewProductRepository cria uma nova instância do ProductRepository com a conexão do banco de dados
func NewProductRepository(connection *sql.DB) ProductRepository {
	return &productRepository{
		connection: connection,
	}
}

// GetProducts busca todos os produtos no banco de dados e retorna uma lista de produtos
func (pr *productRepository) GetProducts(filters model.ProductFilters) ([]model.Product, error) {
	query := `
		SELECT p.products_id, p.products_name, p.products_price, p.categorias_id, c.categorias_name, p.products_data_cadastro, p.products_data_atualizacao, p.products_data_inativacao, p.products_status
		FROM prod.products p
		JOIN prod.categorias c ON p.categorias_id = c.categorias_id`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("p.products_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Categoria_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("p.categorias_id = $%d", argIndex))
		args = append(args, filters.Categoria_Id)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("p.products_status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY p.products_name ASC"

	offset := (filters.Page - 1) * filters.Limit
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, offset)

	rows, err := pr.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productList []model.Product

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Categoria_Id, &product.Categoria_Name, &product.Data_Cadastro, &product.Data_Atualizacao, &product.Data_Inativacao, &product.Status); err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	return productList, nil
}

// GetProductByID busca um produto específico pelo ID no banco de dados
func (pr *productRepository) GetProductByID(id int) (model.Product, error) {
	// Prepara a query SQL para evitar SQL Injection
	query, err := pr.connection.Prepare(`
		SELECT p.products_id, p.products_name, p.products_price, p.categorias_id, c.categorias_name, p.products_data_cadastro, p.products_data_atualizacao, p.products_data_inativacao, p.products_status
		FROM prod.products p
		JOIN prod.categorias c ON p.categorias_id = c.categorias_id
		WHERE p.products_id = $1
	`)
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err) // Log do erro
		return model.Product{}, err                    // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Fecha a consulta após execução

	var product model.Product
	// Executa a query e faz o scan do resultado para o objeto product
	err = query.QueryRow(id).Scan(&product.Id,
		&product.Name,
		&product.Price,
		&product.Categoria_Id,
		&product.Categoria_Name,
		&product.Data_Cadastro,
		&product.Data_Atualizacao,
		&product.Data_Inativacao,
		&product.Status)
	if err != nil {
		// Log de erro caso a consulta falhe
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum produto encontrado com o ID:", id) // Log caso não encontre produto
			return model.Product{}, nil                            // Retorna nil para indicar que o produto não foi encontrado
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return model.Product{}, err                             // Retorna erro caso ocorra outro tipo de falha
	}

	// Log do produto encontrado
	fmt.Println("Produto encontrado:", product)

	return product, nil // Retorna o produto encontrado
}

// CreateProducts insere um novo produto no banco de dados e retorna o ID do produto inserido
func (pr *productRepository) CreateProducts(product model.ProductPost) (int, string, error) {
	var Id int
	var CategoriaName string
	// Prepara a query de inserção para evitar SQL Injection
	query, err := pr.connection.Prepare(`
		INSERT INTO prod.products (products_name, products_price, categorias_id)
		VALUES ($1, $2, $3)
		RETURNING products_id, categorias_id;
	`)
	if err != nil {
		fmt.Println(err)  // Log do erro
		return 0, "", err // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	// Executa a query e escaneia o ID do produto e o ID da categoria associada ao novo produto
	err = query.QueryRow(product.Name, product.Price, product.Categoria_Id).Scan(&Id, &product.Categoria_Id)
	if err != nil {
		fmt.Println(err)  // Log do erro
		return 0, "", err // Retorna erro caso ocorra uma falha na inserção
	}

	// Agora que o produto foi inserido, vamos buscar o nome da categoria com base no categoria_id
	err = pr.connection.QueryRow(`
		SELECT categorias_name
		FROM prod.categorias
		WHERE categorias_id = $1
	`, product.Categoria_Id).Scan(&CategoriaName)

	if err != nil {
		fmt.Println(err)  // Log do erro
		return 0, "", err // Retorna erro caso não consiga recuperar o nome da categoria
	}

	return Id, CategoriaName, nil // Retorna o ID do novo produto e o nome da categoria
}

// UpdateProducts atualiza um produto existente no banco de dados
func (pr *productRepository) UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error) {

	// Prepara a query de atualização
	query, err := pr.connection.Prepare(`
		UPDATE prod.products
		SET products_name = $1, products_price = $2, categorias_id = $3, products_data_atualizacao = CURRENT_TIMESTAMP
		WHERE products_id = $4
		RETURNING products_id, categorias_id;
	`)
	if err != nil {
		fmt.Println(err)                  // Log do erro
		return model.ProductUpdate{}, err // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução
	var categoriaId int
	var CategoriaName string
	// Executa a query e escaneia o ID do produto atualizado
	err = query.QueryRow(product.Name, product.Price, product.Categoria_Id, product.Id).Scan(&product.Id, &categoriaId)
	if err != nil {
		fmt.Println(err)                  // Log do erro
		return model.ProductUpdate{}, err // Retorna erro caso ocorra uma falha na atualização
	}

	err = pr.connection.QueryRow("SELECT categorias_name FROM prod.categorias WHERE categorias_id = $1", categoriaId).Scan(&CategoriaName)
	if err != nil {
		// Tratamento de erro
		fmt.Println("Erro ao buscar o nome da categoria:", err)
		return model.ProductUpdate{}, err
	}
	// Atualiza o nome da categoria no produto
	product.Categoria_Name = CategoriaName
	// Retorna o produto com o ID atualizado
	return product, nil
}

func (pr *productRepository) InactivateProduct(id int) error {
	// Prepara a query de exclusão para evitar SQL Injection
	query, err := pr.connection.Prepare(`
		UPDATE prod.products 
		SET products_data_inativacao = CURRENT_TIMESTAMP, products_status = 'inativo' 
		WHERE products_id = $1
	`)
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

func (pr *productRepository) ActivateProduct(id int) error {
	// Prepara a query de exclusão para evitar SQL Injection
	query, err := pr.connection.Prepare(`
		UPDATE prod.products 
		SET products_data_inativacao = NULL, products_status = 'ativo' 
		WHERE products_id = $1
	`)
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

// ValidateCategory verifica se a categoria existe no banco de dados.
func (pr *productRepository) ValidateCategory(categoriaId int) error {
	return helper.ValidateCategory(pr.connection, categoriaId)
}

func (pr *productRepository) CountProducts(filters model.ProductFilters) (int, error) {
	query := `SELECT COUNT(*) FROM prod.products`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("products_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Categoria_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("categorias_id = $%d", argIndex))
		args = append(args, filters.Categoria_Id)
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("products_status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := pr.connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (pr *productRepository) ValidateProduct(productId int) error {
	return helper.ValidateProduct(pr.connection, productId)
}
