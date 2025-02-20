package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
	"strings"
)

type LstComprasItensRepository interface {
	GetLstComprasItens(filters model.LstCompras_Itens_Filters) ([][]model.LstCompras_Itens, error)
}

type lstComprasItensRepository struct {
	connection *sql.DB
}

func NewLstComprasItensRepository(connection *sql.DB) LstComprasItensRepository {
	return &lstComprasItensRepository{
		connection: connection,
	}
}

func (r *lstComprasItensRepository) GetLstComprasItens(filters model.LstCompras_Itens_Filters) ([][]model.LstCompras_Itens, error) {
	query := `
		SELECT i.lst_compras_itens_id, i.lst_compras_id, i.products_id, p.products_name, 
		       i.lst_compras_itens_quantidade, i.lst_compras_itens_preco, 
		       i.lst_compras_itens_comprado, i.lst_compras_itens_data_cadastro, 
		       i.lst_compras_itens_data_atualizacao
		FROM prod.lst_compras_itens i
		JOIN prod.products p ON i.products_id = p.products_id`

	var conditions []string
	var args []interface{}
	argIndex := 1

	if filters.Product_Name != "" {
		conditions = append(conditions, fmt.Sprintf("p.products_name ILIKE $%d", argIndex))
		args = append(args, "%"+filters.Product_Name+"%")
		argIndex++
	}

	if filters.Product_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("i.products_id = $%d", argIndex))
		args = append(args, filters.Product_Id)
		argIndex++
	}

	if filters.LstCompras_Id > 0 {
		conditions = append(conditions, fmt.Sprintf("i.lst_compras_id = $%d", argIndex))
		args = append(args, filters.LstCompras_Id)
		argIndex++
	}

	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	query += " ORDER BY i.lst_compras_id ASC"

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

	// Mapa para agrupar os itens pelo lst_compras_id
	groupedItems := make(map[int][]model.LstCompras_Itens)

	for rows.Next() {
		var item model.LstCompras_Itens
		if err := rows.Scan(&item.Id, &item.LstCompras_Id, &item.Product_Id, &item.Product_Name, &item.Quantidade, &item.Preco, &item.Item_Check, &item.Data_Cadastro, &item.Data_Atualizacao); err != nil {
			return nil, err
		}
		groupedItems[item.LstCompras_Id] = append(groupedItems[item.LstCompras_Id], item)
	}

	// Converter o mapa para uma slice de slices
	var result [][]model.LstCompras_Itens
	for _, items := range groupedItems {
		result = append(result, items)
	}

	return result, nil
}
