package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

// ProductRepository representa a estrutura do repositório de produtos, armazenando a conexão com o banco de dados
type ProductRepository interface {
	GetProducts(filters model.ProductFilters) ([]model.Product, error)
	GetProductByID(id int) (*model.Product, error)
	CreateProducts(product model.ProductPost) (int, string, error)
	UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error)
	DeleteProduct(id int) error
	ValidateCategory(categoriaId int) error
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
	query := `SELECT products_id, products_name, products_price, categorias_id, categorias_name, products_data_cadastro, products_data_atualizacao FROM prod.products`
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

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY products_name ASC"

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

	rows, err := pr.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productList []model.Product

	for rows.Next() {
		var product model.Product
		if err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.Categoria_Id, &product.Categoria_Name, &product.Data_Cadastro, &product.Data_Atualizacao); err != nil {
			return nil, err
		}
		productList = append(productList, product)
	}

	return productList, nil
}

// GetProductByID busca um produto específico pelo ID no banco de dados
func (pr *productRepository) GetProductByID(id int) (*model.Product, error) {
	// Prepara a query SQL para evitar SQL Injection
	query, err := pr.connection.Prepare("SELECT * FROM prod.products WHERE products_id = $1")
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err) // Log do erro
		return nil, err                                // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Fecha a consulta após execução

	var product model.Product
	// Executa a query e faz o scan do resultado para o objeto product
	err = query.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price, &product.Data_Cadastro, &product.Data_Atualizacao, &product.Categoria_Id, &product.Categoria_Name)
	if err != nil {
		// Log de erro caso a consulta falhe
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum produto encontrado com o ID:", id) // Log caso não encontre produto
			return nil, nil                                        // Retorna nil para indicar que o produto não foi encontrado
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return nil, err                                         // Retorna erro caso ocorra outro tipo de falha
	}

	// Log do produto encontrado
	fmt.Println("Produto encontrado:", product)

	return &product, nil // Retorna o produto encontrado
}

// CreateProducts insere um novo produto no banco de dados e retorna o ID do produto inserido
func (pr *productRepository) CreateProducts(product model.ProductPost) (int, string, error) {
	var Id int
	var CategoriaName string
	// Prepara a query de inserção para evitar SQL Injection
	query, err := pr.connection.Prepare(`
		INSERT INTO prod.products (products_name, products_price, categorias_id, categorias_name)
		SELECT $1, $2, c.categorias_id, c.categorias_name
		FROM prod.categorias c
		WHERE c.categorias_id = $3
		RETURNING products_id, categorias_name;
	`)
	if err != nil {
		fmt.Println(err)  // Log do erro
		return 0, "", err // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	// Executa a query e escaneia o ID e nome da categoria do novo produto inserido
	err = query.QueryRow(product.Name, product.Price, product.Categoria_Id).Scan(&Id, &CategoriaName)
	if err != nil {
		fmt.Println(err)  // Log do erro
		return 0, "", err // Retorna erro caso ocorra uma falha na inserção
	}

	return Id, CategoriaName, nil // Retorna o ID do novo produto e o nome da categoria
}

// UpdateProducts atualiza um produto existente no banco de dados
func (pr *productRepository) UpdateProducts(product model.ProductUpdate) (model.ProductUpdate, error) {

	// Prepara a query de atualização
	query, err := pr.connection.Prepare(`
		UPDATE prod.products 
		SET products_name = $1, products_price = $2, categorias_id = $3, categorias_name = (
			SELECT categorias_name FROM prod.categorias WHERE categorias_id = $3
		), products_data_atualizacao = CURRENT_TIMESTAMP
		WHERE products_id = $4
		RETURNING products_id, categorias_name;
	`)
	if err != nil {
		fmt.Println(err)                  // Log do erro
		return model.ProductUpdate{}, err // Retorna erro caso a preparação da query falhe
	}
	defer query.Close() // Garante que a query será fechada após a execução

	var CategoriaName string
	// Executa a query e escaneia o ID do produto atualizado
	err = query.QueryRow(product.Name, product.Price, product.Categoria_Id, product.Id).Scan(&product.Id, &CategoriaName)
	if err != nil {
		fmt.Println(err)                  // Log do erro
		return model.ProductUpdate{}, err // Retorna erro caso ocorra uma falha na atualização
	}
	// Atualiza o nome da categoria no produto
	product.Categoria_Name = CategoriaName
	// Retorna o produto com o ID atualizado
	return product, nil
}

func (pr *productRepository) DeleteProduct(id int) error {
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

// ValidateCategory verifica se a categoria existe no banco de dados.
func (pr *productRepository) ValidateCategory(categoriaId int) error {
	var count int
	query := `SELECT COUNT(1) FROM prod.categorias WHERE categorias_id = $1`
	err := pr.connection.QueryRow(query, categoriaId).Scan(&count)
	if err != nil {
		return err // Se houver erro ao consultar, retornamos o erro
	}
	if count == 0 {
		return fmt.Errorf("categoria não encontrada") // Se count for 0, categoria não existe
	}
	return nil // Categoria existe
}
