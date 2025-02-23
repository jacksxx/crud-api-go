package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type LstComprasRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetLstCompras(filters model.LstCompras_Filters) ([]model.LstCompras, error)
	GetLstComprasById(id int) (model.LstCompras, error)
	CreateLstCompras(compras model.LstCompras_Post, tx *sql.Tx) (model.LstCompras_Post, error)
	CountLstCompras(filters model.LstCompras_Filters) (int, error)
	TotaisLstCompras(compras model.LstCompras_Post, tx *sql.Tx) (model.LstCompras_Post, error)
	VerificarExistenciaLstCompras(lstComprasId int, tx *sql.Tx) (bool, error)
}

type lstComprasRepository struct {
	connection *sql.DB
}

func NewLstComprasRepository(connection *sql.DB) LstComprasRepository {
	return &lstComprasRepository{
		connection: connection,
	}
}
func (r *lstComprasRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (r *lstComprasRepository) GetLstCompras(filters model.LstCompras_Filters) ([]model.LstCompras, error) {
	query := `
		SELECT lc.lst_compras_id, lc.lst_compras_name, 
		       lc.lst_compras_valor_total, lc.lst_compras_total_itens,
		       lc.lst_compras_status_id, sc.lst_compras_status_name, 
		       lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao
		FROM prod.lst_compras lc
		JOIN prod.lst_compras_status sc ON lc.lst_compras_status_id = sc.lst_compras_status_id`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Nome != "" {
		conditions = append(conditions, fmt.Sprintf("lc.lst_compras_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Nome+"%")
		argIndex++
	}

	if filters.Status_Codigo > 0 {
		conditions = append(conditions, fmt.Sprintf("lc.lst_compras_status_id = $%d", argIndex))
		args = append(args, filters.Status_Codigo)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY lc.lst_compras_id ASC"

	// Paginação
	offset := (filters.Page - 1) * filters.Limit
	if offset < 0 {
		offset = 0
	}
	query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", argIndex, argIndex+1)
	args = append(args, filters.Limit, offset)

	rows, err := r.connection.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var compras []model.LstCompras
	for rows.Next() {
		var compra model.LstCompras
		if err := rows.Scan(&compra.Id, &compra.Nome, &compra.Total, &compra.Qtd_Itens, &compra.Status_Codigo, &compra.Status, &compra.Data_Cadastro, &compra.Data_Atualizacao); err != nil {
			return nil, err
		}
		compras = append(compras, compra)
	}

	return compras, nil

}

func (r *lstComprasRepository) GetLstComprasById(id int) (model.LstCompras, error) {
	query, err := r.connection.Prepare(`
		SELECT lc.lst_compras_id, lc.lst_compras_name, 
		       lc.lst_compras_valor_total, lc.lst_compras_total_itens,
		       lc.lst_compras_status_id, sc.lst_compras_status_name, 
		       lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao
		FROM prod.lst_compras lc
		JOIN prod.lst_compras_status sc ON lc.lst_compras_status_id = sc.lst_compras_status_id
		WHERE lc.lst_compras_id = $1`)
	if err != nil {
		fmt.Println("Error ao preparar consulta:", err)
		return model.LstCompras{}, err
	}
	defer query.Close()
	var compra model.LstCompras

	err = query.QueryRow(id).Scan(&compra.Id, &compra.Nome, &compra.Total, &compra.Qtd_Itens, &compra.Status_Codigo, &compra.Status, &compra.Data_Cadastro, &compra.Data_Atualizacao)
	if err != nil {
		// Log de erro caso a consulta falhe
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum item encontrado com o ID:", id) // Log caso não encontre produto
			return model.LstCompras{}, nil                      // Retorna nil para indicar que o produto não foi encontrado
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return model.LstCompras{}, err                          // Retorna erro caso ocorra outro tipo de falha
	}

	return compra, nil
}

func (r *lstComprasRepository) CreateLstCompras(compras model.LstCompras_Post, tx *sql.Tx) (model.LstCompras_Post, error) {
	var Id int

	query := `
		INSERT INTO prod.lst_compras (lst_compras_name , lst_compras_status_id)
		VALUES ($1, 1)
		RETURNING lst_compras_id
	`
	// Executa a query e escaneia o ID do item inserido
	err := tx.QueryRow(query, compras.Nome).Scan(&Id)
	if err != nil {
		fmt.Println(err)
		return model.LstCompras_Post{}, err
	}
	compras.Id = Id

	return compras, nil
}

func (r *lstComprasRepository) CountLstCompras(filters model.LstCompras_Filters) (int, error) {
	query := `SELECT COUNT(*) FROM prod.lst_compras`
	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Nome != "" {
		conditions = append(conditions, fmt.Sprintf("lst_compras_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Nome+"%")
		argIndex++
	}

	if filters.Status_Codigo > 0 {
		conditions = append(conditions, fmt.Sprintf("lst_compras_status_id = $%d", argIndex))
		args = append(args, filters.Status_Codigo)
		argIndex++
	}

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

func (r *lstComprasRepository) TotaisLstCompras(compras model.LstCompras_Post, tx *sql.Tx) (model.LstCompras_Post, error) {
	// Busca os totais atualizados (status, quantidade de itens, e preço total)

	var status, totalItens int
	var totalPreco float64

	queryTotals, err := r.connection.Prepare(`
		SELECT lst_compras_status_id, lst_compras_total_itens, lst_compras_valor_total
		FROM prod.lst_compras
		WHERE lst_compras_id = $1;
	`)

	if err != nil {
		return model.LstCompras_Post{}, fmt.Errorf("erro ao preparar a consulta para buscar totais atualizados: %v", err)
	}
	err = queryTotals.QueryRow(compras.Id).Scan(&status, &totalItens, &totalPreco)
	if err != nil {
		return model.LstCompras_Post{}, fmt.Errorf("erro ao buscar totais atualizados: %v", err)
	}
	// Atualiza a resposta com os totais calculados
	compras.Status_Codigo = status
	compras.Qtd_Itens = totalItens
	compras.Total = totalPreco

	return compras, nil
}

func (r *lstComprasRepository) VerificarExistenciaLstCompras(lstComprasId int, tx *sql.Tx) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM prod.lst_compras WHERE lst_compras_id = $1)`

	err := tx.QueryRow(query, lstComprasId).Scan(&exists)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil // Retorna false se não encontrar a lista de compras
		}
		return false, fmt.Errorf("erro ao verificar existência da lista de compras: %v", err)
	}

	return exists, nil
}
