package repository

import (
	"crud-api-go/arch/model"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"strings"
)

type CategoriasRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetCategorias(categorias model.CategoriasFilters) ([]model.Categorias, error)
	GetCategoriasById(id int) (model.Categorias, error)
	CreateCategorias(categorias model.CategoriasPost, tx *sql.Tx) (model.CategoriasPost, error)
	UpdateCategoria(categorias model.CategoriasUpdate, tx *sql.Tx) (model.CategoriasUpdate, error)
	InactivateCategoria(id int, tx *sql.Tx) error
	ActivateCategoria(id int, tx *sql.Tx) error
	ValidateCategoryName(nomeCategorias string, categoriaId *int, tx *sql.Tx) error
	ValidateCategory(categoriaId int, tx *sql.Tx) error
	CountCategories(filters model.CategoriasFilters) (int, error)
}

type categoriasRepository struct {
	connection *sql.DB
}

func NewCategoriasRepository(connection *sql.DB) CategoriasRepository {
	return &categoriasRepository{
		connection: connection,
	}
}

func (r *categoriasRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (cr *categoriasRepository) GetCategorias(filters model.CategoriasFilters) ([]model.Categorias, error) {
	query := `SELECT categorias_id, categorias_name, categorias_data_cadastro, categorias_data_atualizacao, categorias_data_inativacao, categorias_status FROM prod.categorias`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("categorias_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("categorias_status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	// Add any additional conditions to the WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	offset := (filters.Page - 1) * filters.Limit
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, offset)

	// Execute the query
	rows, err := cr.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var categoriasList []model.Categorias
	for rows.Next() {
		var categoria model.Categorias
		if err := rows.Scan(&categoria.Id, &categoria.Name, &categoria.Data_cadastro, &categoria.Data_atualizacao, &categoria.Data_Inativacao, &categoria.Status); err != nil {
			return nil, err
		}
		categoriasList = append(categoriasList, categoria)
	}

	return categoriasList, nil
}

func (cr *categoriasRepository) GetCategoriasById(id int) (model.Categorias, error) {
	query, err := cr.connection.Prepare("SELECT * FROM prod.categorias WHERE categorias_id = $1")
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err)
		return model.Categorias{}, err
	}
	defer query.Close()

	var categoria model.Categorias

	err = query.QueryRow(id).Scan(&categoria.Id, &categoria.Name, &categoria.Data_cadastro, &categoria.Data_atualizacao, &categoria.Data_Inativacao, &categoria.Status)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nenhuma categoria encontrada com ID:", id)
			return model.Categorias{}, nil
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return model.Categorias{}, err                          // Retorna erro caso ocorra outro tipo de falha
	}
	fmt.Println("Categoria encontrada:", categoria)

	return categoria, nil
}

func (cr *categoriasRepository) CreateCategorias(categorias model.CategoriasPost, tx *sql.Tx) (model.CategoriasPost, error) {
	var Id int

	query := `
	INSERT INTO prod.categorias (categorias_name) 
	VALUES ($1)
	RETURNING categorias_id`

	err := tx.QueryRow(query, categorias.Name).Scan(&Id)
	if err != nil {
		return model.CategoriasPost{}, fmt.Errorf("erro ao criar categoria: %w", err)
	}
	categorias.Id = Id

	return categorias, nil
}

func (cr *categoriasRepository) UpdateCategoria(categorias model.CategoriasUpdate, tx *sql.Tx) (model.CategoriasUpdate, error) {

	query := `
	UPDATE prod.categorias SET categorias_name = $1, categorias_data_atualizacao = CURRENT_TIMESTAMP 
	WHERE categorias_id = $2
	RETURNING categorias_id`

	err := tx.QueryRow(query, categorias.Name, categorias.Id).Scan(&categorias.Id)
	if err != nil {
		fmt.Println(err)
		return model.CategoriasUpdate{}, err
	}

	return categorias, nil
}

func (cr *categoriasRepository) InactivateCategoria(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.categorias 
		SET categorias_data_inativacao = CURRENT_TIMESTAMP, categorias_status = 'inativo' 
		WHERE categorias_id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (cr *categoriasRepository) ActivateCategoria(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.categorias 
		SET categorias_data_inativacao = NULL, categorias_status = 'ativo' 
		WHERE categorias_id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (cr *categoriasRepository) ValidateCategoryName(nomeCategorias string, categoriaId *int, tx *sql.Tx) error {
	var existingId int

	// Define a query para verificar se já existe uma categoria com o mesmo nome
	// Utiliza ILIKE para fazer a comparação sem diferenciar maiúsculas e minúsculas
	query := `SELECT categorias_id FROM prod.categorias WHERE categorias_name ILIKE $1`
	args := []interface{}{nomeCategorias} // Parâmetro inicial da query (nome da categoria)

	// Se for um update, adiciona uma condição para ignorar a própria categoria
	if categoriaId != nil {
		query += ` AND categorias_id <> $2` // Evita conflito com a própria categoria ao atualizar
		args = append(args, *categoriaId)   // Adiciona o ID da categoria a ser ignorado nos argumentos da query
	}

	// Executa a query e tenta escanear o resultado no existingId
	err := tx.QueryRow(query, args...).Scan(&existingId)

	// Se não houver erro ao escanear, significa que a categoria já existe
	if err == nil {
		return fmt.Errorf("categoria já existe")
	}


	// Se houver erro, mas não for devido a uma categoria existente, retorna nil (nenhum erro encontrado)
	return nil
}

func (cr *categoriasRepository) CountCategories(filters model.CategoriasFilters) (int, error) {
	query := `SELECT COUNT(*) FROM prod.categorias`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("categorias_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("categorias_status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	// Add any additional conditions to the WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := cr.connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func (cr *categoriasRepository) ValidateCategory(categoriaId int, tx *sql.Tx) error {
	return helper.ValidateCategory(tx, categoriaId)
}
