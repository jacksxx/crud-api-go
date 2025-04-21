package repository

import (
	"crud-api-go/arch/model"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"strings"
)

type SubcategoriasRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetSubcategorias(filters model.SubcategoriasFilters) ([]model.Subcategorias, error)
	GetSubcategoriasById(id int) (model.Subcategorias, error)
	CreateSubcategorias(subcategorias model.SubcategoriasPost, tx *sql.Tx) (model.SubcategoriasPost, error)
	UpdateSubcategorias(subcategorias model.SubcategoriasUpdate, tx *sql.Tx) (model.SubcategoriasUpdate, error)
	InactivateSubCategoria(id int, tx *sql.Tx) error
	ActivateSubCategoria(id int, tx *sql.Tx) error
	ValidateSubCategoryName(nomeSubcategorias string, subcategoriaId *int, tx *sql.Tx) error
	CountSubcategories(filters model.SubcategoriasFilters) (int, error)
	ValidateSubcategory(subcategoriaId int, tx *sql.Tx) error
	ValidateCategory(categoriaId int, tx *sql.Tx) error
	CheckStatus(subcategoriaId int, tx *sql.Tx) (string, error)
}

type subcategoriasRepository struct {
	connection *sql.DB
}

func NewSubcategoriasRepository(connection *sql.DB) SubcategoriasRepository {
	return &subcategoriasRepository{
		connection: connection,
	}
}

func (r *subcategoriasRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (r *subcategoriasRepository) GetSubcategorias(filters model.SubcategoriasFilters) ([]model.Subcategorias, error) {
	query := `
		SELECT s.subcategorias_id, s.subcategorias_name, s.categorias_id, c.categorias_name, s.subcategorias_data_cadastro, 
		s.subcategorias_data_atualizacao, s.subcategorias_data_inativacao, s.subcategorias_status 
		FROM prod.subcategorias s
		JOIN prod.categorias c ON s.categorias_id = c.categorias_id
		`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("s.subcategorias_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("s.subcategorias_status = $%d", argIndex))
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
	rows, err := r.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Collect the results
	var subcategoriasList []model.Subcategorias
	for rows.Next() {
		var subcategoria model.Subcategorias
		if err := rows.Scan(&subcategoria.Id, &subcategoria.Name, &subcategoria.CategoriasId, &subcategoria.CategoriasName, &subcategoria.Data_cadastro, &subcategoria.Data_atualizacao, &subcategoria.Data_Inativacao, &subcategoria.Status); err != nil {
			return nil, err
		}
		subcategoriasList = append(subcategoriasList, subcategoria)
	}

	return subcategoriasList, nil
}

func (r *subcategoriasRepository) GetSubcategoriasById(id int) (model.Subcategorias, error) {
	query, err := r.connection.Prepare(`
		SELECT s.subcategorias_id, s.subcategorias_name, s.categorias_id, 
			c.categorias_name, s.subcategorias_data_cadastro, 
			s.subcategorias_data_atualizacao, s.subcategorias_data_inativacao, s.subcategorias_status 
		FROM prod.subcategorias s
		JOIN prod.categorias c ON s.categorias_id = c.categorias_id
		WHERE s.subcategorias_id = $1`)
	if err != nil {
		fmt.Println("Error ao preparar consulta:", err)
		return model.Subcategorias{}, err
	}
	defer query.Close()

	var subcategoria model.Subcategorias

	if err = query.QueryRow(id).Scan(&subcategoria.Id, &subcategoria.Name, &subcategoria.CategoriasId, &subcategoria.CategoriasName, &subcategoria.Data_cadastro, &subcategoria.Data_atualizacao, &subcategoria.Data_Inativacao, &subcategoria.Status); err != nil {
		return model.Subcategorias{}, err
	}

	return subcategoria, nil
}

func (r *subcategoriasRepository) CreateSubcategorias(subcategorias model.SubcategoriasPost, tx *sql.Tx) (model.SubcategoriasPost, error) {
	var Id int
	query, err := tx.Prepare(`
		INSERT INTO prod.subcategorias (subcategorias_name, categorias_id)
		VALUES ($1, $2)
		RETURNING subcategorias_id, subcategorias_name`)
	if err != nil {
		return model.SubcategoriasPost{}, fmt.Errorf("error ao preparar consulta: %w", err)
	}
	defer query.Close()
	err = query.QueryRow(subcategorias.Name, subcategorias.CategoriasId).Scan(&Id, &subcategorias.Name)
	if err != nil {
		return model.SubcategoriasPost{}, fmt.Errorf("erro ao executar insert: %w", err)
	}

	subcategorias.Id = Id
	return subcategorias, nil
}

func (r *subcategoriasRepository) UpdateSubcategorias(subcategorias model.SubcategoriasUpdate, tx *sql.Tx) (model.SubcategoriasUpdate, error) {

	query, err := tx.Prepare(`
		UPDATE prod.subcategorias SET subcategorias_name = $1, categorias_id = $2, subcategorias_data_atualizacao = CURRENT_TIMESTAMP
		WHERE subcategorias_id = $3
		RETURNING subcategorias_id`)
	if err != nil {
		return model.SubcategoriasUpdate{}, fmt.Errorf("error ao preparar consulta: %w", err)
	}
	defer query.Close()

	err = query.QueryRow(subcategorias.Name, subcategorias.CategoriasId, subcategorias.Id).Scan(&subcategorias.Id)
	if err != nil {
		return model.SubcategoriasUpdate{}, fmt.Errorf("erro ao executar insert: %w", err)
	}

	return subcategorias, nil
}

func (r *subcategoriasRepository) InactivateSubCategoria(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.subcategorias 
		SET subcategorias_data_inativacao = CURRENT_TIMESTAMP, subcategorias_status = 'inativo' 
		WHERE subcategorias_id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *subcategoriasRepository) ActivateSubCategoria(id int, tx *sql.Tx) error {
	query := `
		UPDATE prod.subcategorias 
		SET subcategorias_data_inativacao = NULL, subcategorias_status = 'ativo' 
		WHERE subcategorias_id = $1
	`

	_, err := tx.Exec(query, id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (r *subcategoriasRepository) ValidateSubCategoryName(nomeSubcategorias string, subcategoriaId *int, tx *sql.Tx) error {
	var existingId int

	// Define a query para verificar se já existe uma categoria com o mesmo nome
	// Utiliza ILIKE para fazer a comparação sem diferenciar maiúsculas e minúsculas
	query := `SELECT subcategorias_id FROM prod.subcategorias WHERE subcategorias_name ILIKE $1`
	args := []interface{}{nomeSubcategorias} // Parâmetro inicial da query (nome da categoria)

	// Se for um update, adiciona uma condição para ignorar a própria categoria
	if subcategoriaId != nil {
		query += ` AND subcategorias_id <> $2` // Evita conflito com a própria categoria ao atualizar
		args = append(args, *subcategoriaId)   // Adiciona o ID da categoria a ser ignorado nos argumentos da query
	}

	// Executa a query e tenta escanear o resultado no existingId
	err := tx.QueryRow(query, args...).Scan(&existingId)

	// Se não houver erro ao escanear, significa que a categoria já existe
	if err == nil {
		return fmt.Errorf("subcategoria já existe")
	}

	// Se houver erro, mas não for devido a uma categoria existente, retorna nil (nenhum erro encontrado)
	return nil
}

//Helpers and Utils

func (r *subcategoriasRepository) CountSubcategories(filters model.SubcategoriasFilters) (int, error) {
	query := `SELECT COUNT(*) FROM prod.subcategorias`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("subcategorias_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	if filters.Status != "" {
		conditions = append(conditions, fmt.Sprintf("subcategorias_status = $%d", argIndex))
		args = append(args, filters.Status)
		argIndex++
	}

	// Add any additional conditions to the WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	var count int
	err := r.connection.QueryRow(query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func (r *subcategoriasRepository) ValidateSubcategory(subcategoriaId int, tx *sql.Tx) error {
	return helper.ValidateSubCategory(tx, subcategoriaId)
}

// ValidateCategory verifica se a categoria existe no banco de dados.
func (r *subcategoriasRepository) ValidateCategory(categoriaId int, tx *sql.Tx) error {
	return helper.ValidateCategory(tx, categoriaId)
}

func (r *subcategoriasRepository) CheckStatus(subcategoriaId int, tx *sql.Tx) (string, error) {
	var status string

	query := `
		SELECT subcategorias_status 
		FROM prod.subcategorias 
		WHERE subcategorias_id = $1
	`

	err := tx.QueryRow(query, subcategoriaId).Scan(&status)
	if err != nil {

		return "", fmt.Errorf("erro ao verificar status da subcategoria: %v", err)
	}

	return status, nil
}
