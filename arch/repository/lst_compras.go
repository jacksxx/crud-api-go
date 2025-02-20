package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type LstComprasRepository interface {
	GetLstCompras(filters model.LstCompras_Filters) ([]model.LstCompras, error)
	CountLstCompras(filters model.LstCompras_Filters) (int, error)
}

type lstComprasRepository struct {
	connection *sql.DB
}

func NewLstComprasRepository(connection *sql.DB) LstComprasRepository {
	return &lstComprasRepository{
		connection: connection,
	}
}

func (r *lstComprasRepository) GetLstCompras(filters model.LstCompras_Filters) ([]model.LstCompras, error) {
	query := `
		SELECT lc.lst_compras_id, lc.lst_compras_name, 
		       SUM(i.lst_compras_itens_preco * i.lst_compras_itens_quantidade) AS lst_compras_valor_total,
		       COUNT(i.products_id) AS lst_compras_total_itens,
		       lc.lst_compras_status_id, sc.lst_compras_status_name, 
		       lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao
		FROM prod.lst_compras lc
		JOIN prod.lst_compras_status sc ON lc.lst_compras_status_id = sc.lst_compras_status_id
		LEFT JOIN prod.lst_compras_itens i ON lc.lst_compras_id = i.lst_compras_id`

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
	query += " GROUP BY lc.lst_compras_id, lc.lst_compras_name, lc.lst_compras_status_id, sc.lst_compras_status_name, lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao"
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
