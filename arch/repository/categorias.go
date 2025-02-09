package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type CategoriasRepository interface {
	GetCategorias(categorias model.CategoriasFilters) ([]model.Categorias, error)
	GetCategoriasById(id int) (*model.Categorias, error)
	CreateCategorias(categorias model.CategoriasPost) (int, error)
	UpdateCategoria(categorias model.CategoriasUpdate) (model.CategoriasUpdate, error)
	DeleteCategoria(id int) error
	ValidateCategoryName(nomeCategorias string) error
}

type categoriasRepository struct {
	connection *sql.DB
}

func NewCategoriasRepository(connection *sql.DB) CategoriasRepository {
	return &categoriasRepository{
		connection: connection,
	}
}

// TODO: CORRIGIR FILTRO POR NOME
func (cr *categoriasRepository) GetCategorias(filters model.CategoriasFilters) ([]model.Categorias, error) {
	query := `SELECT categorias_id, categorias_name, categorias_data_cadastro, categorias_data_atualizacao FROM prod.categorias`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Name != "" {
		conditions = append(conditions, fmt.Sprintf("categorias_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Name+"%")
		argIndex++
	}

	// Add any additional conditions to the WHERE clause
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	// Add ORDER BY clause
	query += " ORDER BY categorias_name ASC "

	// Set pagination
	limit := max(filters.Limit, 1)
	page := max(filters.Page, 1)
	offset := (page - 1) * limit

	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, limit, offset)

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
		if err := rows.Scan(&categoria.Id, &categoria.Name, &categoria.Data_cadastro, &categoria.Data_atualizacao); err != nil {
			return nil, err
		}
		categoriasList = append(categoriasList, categoria)
	}

	return categoriasList, nil
}

func (cr *categoriasRepository) GetCategoriasById(id int) (*model.Categorias, error) {
	query, err := cr.connection.Prepare("SELECT * FROM prod.categorias WHERE categorias_id = $1")
	if err != nil {
		fmt.Println("Erro ao preparar consulta:", err)
		return nil, err
	}
	defer query.Close()

	var categoria model.Categorias

	err = query.QueryRow(id).Scan(&categoria.Id, &categoria.Name, &categoria.Data_cadastro, &categoria.Data_atualizacao)

	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("Nenhuma categoria encontrada com ID:", id)
			return nil, nil
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return nil, err                                         // Retorna erro caso ocorra outro tipo de falha
	}
	fmt.Println("Categoria encontrada:", categoria)

	return &categoria, nil
}

func (cr *categoriasRepository) CreateCategorias(categorias model.CategoriasPost) (int, error) {
	var Id int

	query, err := cr.connection.Prepare(`
	INSERT INTO prod.categorias (categorias_name) 
	VALUES ($1)
	RETURNING categorias_id`)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer query.Close()

	err = query.QueryRow(categorias.Name).Scan(&Id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return Id, nil
}

func (cr *categoriasRepository) UpdateCategoria(categorias model.CategoriasUpdate) (model.CategoriasUpdate, error) {

	query, err := cr.connection.Prepare(`
	UPDATE prod.categorias SET categorias_name = $1, categorias_data_atualizacao = CURRENT_TIMESTAMP 
	WHERE categorias_id = $2
	RETURNING categorias_id`)
	if err != nil {
		fmt.Println(err)
		return model.CategoriasUpdate{}, err
	}
	defer query.Close()

	err = query.QueryRow(categorias.Name, categorias.Id).Scan(&categorias.Id)
	if err != nil {
		fmt.Println(err)
		return model.CategoriasUpdate{}, err
	}

	return categorias, nil
}

func (cr *categoriasRepository) DeleteCategoria(id int) error {
	query, err := cr.connection.Prepare("DELETE FROM prod.categorias WHERE categorias_id = $1")
	if err != nil {
		fmt.Println(err)
		return err
	}
	defer query.Close()

	_, err = query.Exec(id)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func (cs *categoriasRepository) ValidateCategoryName(nomeCategorias string) error {
	var existingId int

	// Verifica se a categoria já existe ignorando maiúsculas e minúsculas
	err := cs.connection.QueryRow(`
		SELECT categorias_id FROM prod.categorias WHERE categorias_name ILIKE $1
	`, nomeCategorias).Scan(&existingId)

	if err == nil {
		return fmt.Errorf("categoria já existe")
	}

	return nil

}
