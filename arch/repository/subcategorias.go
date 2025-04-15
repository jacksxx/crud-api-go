package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type SubcategoriasRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetSubcategorias(filters model.SubcategoriasFilters) ([]model.Subcategorias, error)
	CountSubcategories(filters model.SubcategoriasFilters) (int, error)
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
