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
	BeginTransaction() (*sql.Tx, error)
	GetProducts(filters model.ProductFilters) ([]model.Product, error)
	GetProductByID(id int) (model.Product, error)
	CreateProducts(product model.ProductPost, tx *sql.Tx) (model.ProductPost, error)
	UpdateProducts(product model.ProductUpdate, tx *sql.Tx) (model.ProductUpdate, error)
	InactivateProduct(id int, tx *sql.Tx) error
	ActivateProduct(id int, tx *sql.Tx) error
	ValidateCategory(categoriaId int, tx *sql.Tx) error
	CountProducts(filters model.ProductFilters) (int, error)
	ValidateProduct(productId int, tx *sql.Tx) error
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

func (r *productRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

// GetProducts busca todos os produtos no banco de dados e retorna uma lista de produtos
func (pr *productRepository) GetProducts(filters model.ProductFilters) ([]model.Product, error) {
	query := `
		SELECT p.products_id, p.products_name, p.products_price, p.categorias_id, c.categorias_name, p.unidade_id, u.unidade_descricao, u.unidade_abreviacao, p.products_data_cadastro, p.products_data_atualizacao, p.products_data_inativacao, p.products_status
		FROM prod.products p
		JOIN prod.categorias c ON p.categorias_id = c.categorias_id
		JOIN prod.unidades u ON p.unidade_id = u.unidade_id`
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

	if filters.Unidade_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("p.unidade_id = $%d", argIndex))
		args = append(args, filters.Unidade_Id)
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
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Categoria_Id, &product.Categoria_Name, &product.Unidade_Id, &product.Unidade_Descricao, &product.Unidade_Abreviacao, &product.Data_Cadastro, &product.Data_Atualizacao, &product.Data_Inativacao, &product.Status); err != nil {
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
		SELECT p.products_id, p.products_name, p.products_price, p.categorias_id, c.categorias_name, p.unidade_id, u.unidade_descricao, u.unidade_abreviacao, p.products_data_cadastro, p.products_data_atualizacao, p.products_data_inativacao, p.products_status
		FROM prod.products p
		JOIN prod.categorias c ON p.categorias_id = c.categorias_id
		JOIN prod.unidades u ON p.unidade_id = u.unidade_id
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
		&product.Unidade_Id,
		&product.Unidade_Descricao,
		&product.Unidade_Abreviacao,
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
func (pr *productRepository) CreateProducts(product model.ProductPost, tx *sql.Tx) (model.ProductPost, error) {
	var Id int
	var CategoriaName string
	var UnidadeDescricao string

	// Prepara a query de inserção para evitar SQL Injection
	query := `
		INSERT INTO prod.products (products_name, products_price, categorias_id, unidade_id)
		VALUES ($1, $2, $3, $4)
		RETURNING products_id, categorias_id, unidade_id;
	`

	// Executa a query e escaneia o ID do produto e os IDs da categoria e unidade associadas
	err := tx.QueryRow(query, product.Name, product.Price, product.Categoria_Id, product.Unidade_Id).
		Scan(&Id, &product.Categoria_Id, &product.Unidade_Id)
	if err != nil {
		fmt.Println(err)                // Log do erro
		return model.ProductPost{}, err // Retorna erro caso ocorra uma falha na inserção
	}

	// Agora que o produto foi inserido, vamos buscar o nome da categoria com base no categoria_id
	err = tx.QueryRow(`
		SELECT categorias_name
		FROM prod.categorias
		WHERE categorias_id = $1
	`, product.Categoria_Id).Scan(&CategoriaName)

	if err != nil {
		fmt.Println(err)                // Log do erro
		return model.ProductPost{}, err // Retorna erro caso não consiga recuperar o nome da categoria
	}

	// Agora, buscar a descrição da unidade com base no unidade_id
	err = tx.QueryRow(`
		SELECT unidade_descricao
		FROM prod.unidades
		WHERE unidade_id = $1
	`, product.Unidade_Id).Scan(&UnidadeDescricao)

	if err != nil {
		fmt.Println(err)                // Log do erro
		return model.ProductPost{}, err // Retorna erro caso não consiga recuperar a descrição da unidade
	}

	// Atribui o nome da categoria e a descrição da unidade aos campos do produto
	product.Id = Id
	product.Categoria_Name = CategoriaName
	product.Unidade_Descricao = UnidadeDescricao

	return product, nil // Retorna o produto com o ID, nome da categoria e descrição da unidade
}

// UpdateProducts atualiza um produto existente no banco de dados
func (pr *productRepository) UpdateProducts(product model.ProductUpdate, tx *sql.Tx) (model.ProductUpdate, error) {

	// Prepara a query de atualização
	query := `
		UPDATE prod.products
		SET products_name = $1, products_price = $2, categorias_id = $3, products_data_atualizacao = CURRENT_TIMESTAMP
		WHERE products_id = $4
		RETURNING products_id, categorias_id, unidade_id;
	`
	// Executa a query de atualização e escaneia os valores retornados
	var categoriaId, unidadeId int
	err := tx.QueryRow(query, product.Name, product.Price, product.Categoria_Id, product.Id).
		Scan(&product.Id, &categoriaId, &unidadeId)
	if err != nil {
		fmt.Println("Erro ao executar a query de atualização:", err)
		return model.ProductUpdate{}, err // Retorna erro caso ocorra uma falha na atualização
	}

	// Busca o nome da categoria com base no categoriaId
	var categoriaName string
	err = tx.QueryRow("SELECT categorias_name FROM prod.categorias WHERE categorias_id = $1", categoriaId).Scan(&categoriaName)
	if err != nil {
		fmt.Println("Erro ao buscar o nome da categoria:", err)
		return model.ProductUpdate{}, err // Retorna erro caso não consiga recuperar o nome da categoria
	}

	// Busca a descrição da unidade com base no unidadeId
	var unidadeDescricao string
	err = tx.QueryRow("SELECT unidade_descricao FROM prod.unidades WHERE unidade_id = $1", unidadeId).Scan(&unidadeDescricao)
	if err != nil {
		fmt.Println("Erro ao buscar a descrição da unidade:", err)
		return model.ProductUpdate{}, err // Retorna erro caso não consiga recuperar a descrição da unidade
	}

	// Atualiza os campos do produto com os valores obtidos
	product.Categoria_Name = categoriaName
	product.Unidade_Descricao = unidadeDescricao

	// Retorna o produto com as informações atualizadas
	return product, nil
}

func (pr *productRepository) InactivateProduct(id int, tx *sql.Tx) error {
	// Prepara a query de exclusão para evitar SQL Injection
	query := `
		UPDATE prod.products 
		SET products_data_inativacao = CURRENT_TIMESTAMP, products_status = 'inativo' 
		WHERE products_id = $1
	` // Garante que a query será fechada após a execução

	// Executa a query de exclusão
	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err) // Log do erro
		return err       // Retorna erro caso ocorra uma falha na exclusão
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}

func (pr *productRepository) ActivateProduct(id int, tx *sql.Tx) error {
	// Prepara a query de exclusão para evitar SQL Injection
	query := `
		UPDATE prod.products 
		SET products_data_inativacao = NULL, products_status = 'ativo' 
		WHERE products_id = $1
	` // Garante que a query será fechada após a execução

	// Executa a query de exclusão
	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err) // Log do erro
		return err       // Retorna erro caso ocorra uma falha na exclusão
	}
	return nil // Retorna nil para indicar sucesso na exclusão
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func (pr *productRepository) ValidateCategory(categoriaId int, tx *sql.Tx) error {
	return helper.ValidateCategory(tx, categoriaId)
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

	if filters.Unidade_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("unidade_id = $%d", argIndex))
		args = append(args, filters.Unidade_Id)
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

func (pr *productRepository) ValidateProduct(productId int, tx *sql.Tx) error {
	return helper.ValidateProduct(tx, productId)
}
