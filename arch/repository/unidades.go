package repository

import (
	"crud-api-go/arch/model"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"strings"
)

type UnidadesRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetUnidades(filters model.UnidadesFilters) ([]model.Unidades, error)
	GetUnidadesById(id int) (model.Unidades, error)
	CreateUnidades(unidade model.UnidadesPost, tx *sql.Tx) (model.UnidadesPost, error)
	UpdateUnidades(unidade model.UnidadesUpdate, tx *sql.Tx) (model.UnidadesUpdate, error)
	ValidateUnitName(unidadeNome string, unidadeId *int, tx *sql.Tx) error
	ValidateUnitAbrev(unidadeAbrev string, unidadeId *int, tx *sql.Tx) error
	CountUnits(filters model.UnidadesFilters) (int, error)
	ValidateUnit(categoriaId int, tx *sql.Tx) error
}

type unidadesRepository struct {
	connection *sql.DB
}

func NewUnidadesRepository(connection *sql.DB) UnidadesRepository {
	return &unidadesRepository{
		connection: connection,
	}
}

func (r *unidadesRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (r *unidadesRepository) GetUnidades(filters model.UnidadesFilters) ([]model.Unidades, error) {
	query := `
		SELECT unidade_id, unidade_descricao, unidade_abreviacao, unidade_data_cadastro, unidade_data_atualizacao 
		FROM prod.unidades
	`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Descricao != "" {
		conditions = append(conditions, fmt.Sprintf(" unidade_descricao ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Descricao+"%")
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
	var unidades []model.Unidades
	for rows.Next() {
		var unidade model.Unidades
		if err := rows.Scan(&unidade.Id, &unidade.Descricao, &unidade.Abreviacao, &unidade.Data_cadastro, &unidade.Data_atualizacao); err != nil {
			return nil, err
		}
		unidades = append(unidades, unidade)
	}

	return unidades, nil
}

func (r *unidadesRepository) GetUnidadesById(id int) (model.Unidades, error) {
	query, err := r.connection.Prepare("SELECT * FROM prod.unidades WHERE unidade_id = $1")
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err)
		return model.Unidades{}, err
	}
	defer query.Close()

	var unidade model.Unidades

	err = query.QueryRow(id).Scan(&unidade.Id, &unidade.Descricao, &unidade.Abreviacao, &unidade.Data_cadastro, &unidade.Data_atualizacao)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nenhuma unidade encontrada com ID:", id)
			return model.Unidades{}, nil
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return model.Unidades{}, err                            // Retorna erro caso ocorra outro tipo de falha
	}
	fmt.Println("Unidade encontrada:", unidade)

	return unidade, nil
}

func (r *unidadesRepository) CreateUnidades(unidade model.UnidadesPost, tx *sql.Tx) (model.UnidadesPost, error) {
	var Id int

	query := `
		INSERT INTO prod.unidades (unidade_descricao , unidade_abreviacao)
		VALUES ($1, $2)
		RETURNING unidade_id
	`
	// Executa a query e escaneia o ID do item inserido
	err := tx.QueryRow(query, unidade.Descricao, unidade.Abreviacao).Scan(&Id)
	if err != nil {
		fmt.Println(err)
		return model.UnidadesPost{}, err
	}
	unidade.Id = Id

	return unidade, nil
}

func (r *unidadesRepository) UpdateUnidades(unidade model.UnidadesUpdate, tx *sql.Tx) (model.UnidadesUpdate, error) {

	query := `
		UPDATE prod.unidades SET unidade_descricao = $1, unidade_abreviacao = $2, unidade_data_atualizacao = CURRENT_TIMESTAMP
		WHERE unidade_id = $3
		RETURNING unidade_id
	`
	// Executa a query e escaneia o ID do item inserido
	err := tx.QueryRow(query, unidade.Descricao, unidade.Abreviacao, unidade.Id).Scan(&unidade.Id)
	if err != nil {
		fmt.Println(err)
		return model.UnidadesUpdate{}, err
	}

	return unidade, nil
}

func (r *unidadesRepository) ValidateUnitName(unidadeNome string, unidadeId *int, tx *sql.Tx) error {
	var existingId int

	// Define a query para verificar se já existe uma categoria com o mesmo nome
	// Utiliza ILIKE para fazer a comparação sem diferenciar maiúsculas e minúsculas
	query := `SELECT unidade_id FROM prod.unidades WHERE unidade_descricao ILIKE $1`
	args := []interface{}{unidadeNome} // Parâmetro inicial da query (nome da categoria)

	// Se for um update, adiciona uma condição para ignorar a própria categoria
	if unidadeId != nil {
		query += ` AND unidade_id <> $2` // Evita conflito com a própria categoria ao atualizar
		args = append(args, *unidadeId)  // Adiciona o ID da categoria a ser ignorado nos argumentos da query
	}

	// Executa a query e tenta escanear o resultado no existingId
	err := tx.QueryRow(query, args...).Scan(&existingId)

	// Se não houver erro ao escanear, significa que a categoria já existe
	if err == nil {
		return fmt.Errorf("unidade já existe")
	}

	// Se houver erro, mas não for devido a uma categoria existente, retorna nil (nenhum erro encontrado)
	return nil
}

func (r *unidadesRepository) ValidateUnitAbrev(unidadeAbrev string, unidadeId *int, tx *sql.Tx) error {
	var existingId int

	// Define a query para verificar se já existe uma categoria com o mesmo nome
	// Utiliza ILIKE para fazer a comparação sem diferenciar maiúsculas e minúsculas
	query := `SELECT unidade_id FROM prod.unidades WHERE unidade_abreviacao ILIKE $1`
	args := []interface{}{unidadeAbrev} // Parâmetro inicial da query (nome da categoria)

	// Se for um update, adiciona uma condição para ignorar a própria categoria
	if unidadeId != nil {
		query += ` AND unidade_id <> $2` // Evita conflito com a própria categoria ao atualizar
		args = append(args, *unidadeId)  // Adiciona o ID da categoria a ser ignorado nos argumentos da query
	}

	// Executa a query e tenta escanear o resultado no existingId
	err := tx.QueryRow(query, args...).Scan(&existingId)

	// Se não houver erro ao escanear, significa que a categoria já existe
	if err == nil {
		return fmt.Errorf("abreviação já existe")
	}

	// Se houver erro, mas não for devido a uma categoria existente, retorna nil (nenhum erro encontrado)
	return nil
}

func (r *unidadesRepository) CountUnits(filters model.UnidadesFilters) (int, error) {
	query := `SELECT COUNT(*) FROM prod.unidades`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Descricao != "" {
		conditions = append(conditions, fmt.Sprintf(" unidade_descricao ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Descricao+"%")
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

func (r *unidadesRepository) ValidateUnit(categoriaId int, tx *sql.Tx) error {
	return helper.ValidateUnit(tx, categoriaId)
}
