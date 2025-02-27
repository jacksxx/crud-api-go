package repository

import (
	"crud-api-go/arch/model"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"strings"
)

type LstComprasItensRepository interface {
	BeginTransaction() (*sql.Tx, error)
	GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int][]model.LstCompras_Itens, error)
	GetLstComprasItensById(id int) (model.LstCompras_Itens, error)
	CreateLstComprasItem(item model.LstCompras_Itens_Post, tx *sql.Tx) (model.LstCompras_Itens_Post, error)
	ValidateProduct(productId int, tx *sql.Tx) error
	UpdateLstComprasTotals(lstComprasId int, tx *sql.Tx) error
	CheckLstComprasExists(lstComprasId int, tx *sql.Tx) (bool, error)
}

type lstComprasItensRepository struct {
	connection *sql.DB
}

func NewLstComprasItensRepository(connection *sql.DB) LstComprasItensRepository {
	return &lstComprasItensRepository{
		connection: connection,
	}
}

func (r *lstComprasItensRepository) BeginTransaction() (*sql.Tx, error) {
	return r.connection.Begin()
}

func (r *lstComprasItensRepository) GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int][]model.LstCompras_Itens, error) {
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

	// Criar um mapa para agrupar os itens por lst_compras_id
	groupedItems := make(map[int][]model.LstCompras_Itens)

	for rows.Next() {
		var item model.LstCompras_Itens
		if err := rows.Scan(&item.Id, &item.LstCompras_Id, &item.Product_Id, &item.Product_Name, &item.Quantidade, &item.Preco, &item.Item_Check, &item.Data_Cadastro, &item.Data_Atualizacao); err != nil {
			return nil, err
		}
		groupedItems[item.LstCompras_Id] = append(groupedItems[item.LstCompras_Id], item)
	}

	return groupedItems, nil
}

func (r *lstComprasItensRepository) GetLstComprasItensById(id int) (model.LstCompras_Itens, error) {
	query, err := r.connection.Prepare(`
		SELECT i.lst_compras_itens_id, i.lst_compras_id, i.products_id, p.products_name, 
		       i.lst_compras_itens_quantidade, i.lst_compras_itens_preco, 
		       i.lst_compras_itens_comprado, i.lst_compras_itens_data_cadastro, 
		       i.lst_compras_itens_data_atualizacao
		FROM prod.lst_compras_itens i
		JOIN prod.products p ON i.products_id = p.products_id
		WHERE i.lst_compras_itens_id = $1
	`)
	if err != nil {
		fmt.Println("Error ao preparar consulta:", err)
		return model.LstCompras_Itens{}, err
	}

	defer query.Close()

	var item model.LstCompras_Itens

	err = query.QueryRow(id).Scan(&item.Id, &item.LstCompras_Id, &item.Product_Id, &item.Product_Name, &item.Quantidade, &item.Preco, &item.Item_Check, &item.Data_Cadastro, &item.Data_Atualizacao)
	if err != nil {
		// Log de erro caso a consulta falhe
		if err == sql.ErrNoRows {
			fmt.Println("Nenhum item encontrado com o ID:", id) // Log caso não encontre produto
			return model.LstCompras_Itens{}, nil                // Retorna nil para indicar que o produto não foi encontrado
		}
		fmt.Println("Erro na consulta ao banco de dados:", err) // Log do erro
		return model.LstCompras_Itens{}, err                    // Retorna erro caso ocorra outro tipo de falha
	}

	return item, nil
}

func (r *lstComprasItensRepository) CreateLstComprasItem(item model.LstCompras_Itens_Post, tx *sql.Tx) (model.LstCompras_Itens_Post, error) {
	if item.LstCompras_Id == 0 {
		return model.LstCompras_Itens_Post{}, fmt.Errorf("erro: LstCompras_Id inválido (0) ao inserir item")
	}

	var Id int
	var ProductName string

	// Inserindo o item na tabela lst_compras_itens e fazendo a junção com a tabela de produtos para pegar o nome do produto
	query := `
		INSERT INTO prod.lst_compras_itens (lst_compras_id, products_id, lst_compras_itens_quantidade, lst_compras_itens_preco)
		VALUES ($1, $2, $3, $4)
		RETURNING lst_compras_itens_id
	`
	// Executando a inserção do item
	err := tx.QueryRow(query, item.LstCompras_Id, item.Product_Id, item.Quantidade, item.Preco).Scan(&Id)
	if err != nil {
		return model.LstCompras_Itens_Post{}, err
	}

	// Agora, realizando a junção para pegar o nome do produto usando o ID do produto inserido
	queryGetProduct := `
		SELECT p.products_name
		FROM prod.products p
		WHERE p.products_id = $1
	`

	// Executando a consulta para obter o nome do produto
	err = tx.QueryRow(queryGetProduct, item.Product_Id).Scan(&ProductName)
	if err != nil {
		return model.LstCompras_Itens_Post{}, err
	}

	item.Id = Id
	item.Product_Name = ProductName

	return item, nil
}

func (r *lstComprasItensRepository) ValidateProduct(productId int, tx *sql.Tx) error {
	return helper.ValidateProduct(tx, productId)
}

func (r *lstComprasItensRepository) UpdateLstComprasTotals(lstComprasId int, tx *sql.Tx) error {
	query, err := tx.Prepare(`
		UPDATE prod.lst_compras
		SET lst_compras_valor_total = (
			SELECT COALESCE(SUM(lst_compras_itens_quantidade * lst_compras_itens_preco), 0)
			FROM prod.lst_compras_itens
			WHERE lst_compras_id = $1
		),
		lst_compras_total_itens = (
			SELECT COALESCE(SUM(lst_compras_itens_quantidade), 0)
			FROM prod.lst_compras_itens
			WHERE lst_compras_id = $1
		)
		WHERE lst_compras_id = $1;
	`)
	if err != nil {
		return err
	}
	defer query.Close()
	
	_, err = query.Exec(lstComprasId)
	if err != nil {
		fmt.Println("Erro ao atualizar totais da lista de compras:", err)
		return err
	}

	return nil
}

// TODO: breno
// Nova função para verificar se LstCompras_Id existe antes de inserir
func (r *lstComprasItensRepository) CheckLstComprasExists(lstComprasId int, tx *sql.Tx) (bool, error) {
	var exists bool
	query := `
		SELECT EXISTS(SELECT 1 FROM prod.lst_compras WHERE lst_compras_id = $1)
	`
	err := tx.QueryRow(query, lstComprasId).Scan(&exists)
	return exists, err
}
