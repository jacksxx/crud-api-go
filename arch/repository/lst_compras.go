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
	UpdateLstCompras(compras model.LstCompras_Update, tx *sql.Tx) (model.LstCompras_Update, error)
	FinishLstCompras(compras model.LstCompras_Finish, tx *sql.Tx) (model.LstCompras_Finish, error)
	CountLstCompras(filters model.LstCompras_Filters) (int, error)
	TotaisLstCompras(compras model.LstCompras_Post, tx *sql.Tx) (model.LstCompras_Post, error)
	VerificarExistenciaLstCompras(lstComprasId int, tx *sql.Tx) (bool, error)
	CancelLstCompras(compras model.LstCompras_Cancel, tx *sql.Tx) error
	BuscarStatusLstCompras(lstComprasId int, tx *sql.Tx) (int, error)
	VerificarStatusLstCompras(lstComprasId int, tx *sql.Tx) (int, error)
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
		       lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao,
			   lc.lst_compras_data_cancelamento, lc.lst_compras_data_finalizacao
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
		if err := rows.Scan(&compra.Id, &compra.Nome, &compra.Total, &compra.Qtd_Itens, &compra.Status_Codigo, &compra.Status, &compra.Data_Cadastro, &compra.Data_Atualizacao, &compra.Data_Cancelamento, &compra.Data_Finalizacao); err != nil {
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
		       lc.lst_compras_data_cadastro, lc.lst_compras_data_atualizacao,
			   lc.lst_compras_data_cancelamento, lc.lst_compras_data_finalizacao
		FROM prod.lst_compras lc
		JOIN prod.lst_compras_status sc ON lc.lst_compras_status_id = sc.lst_compras_status_id
		WHERE lc.lst_compras_id = $1`)
	if err != nil {
		fmt.Println("Error ao preparar consulta:", err)
		return model.LstCompras{}, err
	}
	defer query.Close()
	var compra model.LstCompras

	err = query.QueryRow(id).Scan(&compra.Id, &compra.Nome, &compra.Total, &compra.Qtd_Itens, &compra.Status_Codigo, &compra.Status, &compra.Data_Cadastro, &compra.Data_Atualizacao, &compra.Data_Cancelamento, &compra.Data_Finalizacao)
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
	var status, totalItens int
	var totalPreco float64
	var status_string string
	query := `
		INSERT INTO prod.lst_compras (lst_compras_name)
		VALUES ($1)
		RETURNING lst_compras_id, lst_compras_status_id, lst_compras_total_itens, lst_compras_valor_total
	`
	// Executa a query e escaneia o ID do item inserido
	err := tx.QueryRow(query, compras.Nome).Scan(&Id, &status, &totalItens, &totalPreco)
	if err != nil {
		fmt.Println(err)
		return model.LstCompras_Post{}, err
	}
	queryA := `
		SELECT lst_compras_status_name
		FROM prod.lst_compras_status
		WHERE lst_compras_status_id = $1
	`
	// Executando a consulta para obter o nome do produto
	err = tx.QueryRow(queryA, status).Scan(&status_string)
	if err != nil {
		return model.LstCompras_Post{}, err
	}

	compras.Status = status_string
	compras.Id = Id
	compras.Status_Codigo = status
	compras.Qtd_Itens = totalItens
	compras.Total = totalPreco

	return compras, nil
}

func (r *lstComprasRepository) UpdateLstCompras(compras model.LstCompras_Update, tx *sql.Tx) (model.LstCompras_Update, error) {
	var status_string string
	// if compras.Status_Codigo != 1 {
	// 	return model.LstCompras_Update{}, fmt.Errorf("A Lista de Compra não se encontra em andamento")
	// }
	query := `
		UPDATE prod.lst_compras
		SET lst_compras_name = $1, lst_compras_data_atualizacao = CURRENT_TIMESTAMP
		WHERE lst_compras_id = $2
		RETURNING lst_compras_status_id, lst_compras_total_itens, lst_compras_valor_total
	`
	// Executa a atualização, mas não altera o total de itens e o valor total
	err := tx.QueryRow(query, compras.Nome, compras.Id).Scan(&compras.Status_Codigo, &compras.Qtd_Itens, &compras.Total)
	if err != nil {
		return model.LstCompras_Update{}, fmt.Errorf("erro ao atualizar lista de compras: %w", err)
	}

	queryA := `
		SELECT lst_compras_status_name
		FROM prod.lst_compras_status
		WHERE lst_compras_status_id = $1
	`
	// Obtendo o nome do status atualizado
	err = tx.QueryRow(queryA, compras.Status_Codigo).Scan(&status_string)
	if err != nil {
		return model.LstCompras_Update{}, err
	}

	compras.Status = status_string

	return compras, nil
}

func (r *lstComprasRepository) FinishLstCompras(compras model.LstCompras_Finish, tx *sql.Tx) (model.LstCompras_Finish, error) {
	var status_string string

	query := `
		UPDATE prod.lst_compras
		SET lst_compras_status_id = 2, lst_compras_data_atualizacao = CURRENT_TIMESTAMP
		WHERE lst_compras_id = $1
		RETURNING lst_compras_status_id, lst_compras_total_itens, lst_compras_valor_total, lst_compras_name
	`
	// Executa a atualização, mas não altera o total de itens e o valor total
	err := tx.QueryRow(query, compras.Id).Scan(&compras.Status_Codigo, &compras.Qtd_Itens, &compras.Total, &compras.Nome)
	if err != nil {
		return model.LstCompras_Finish{}, fmt.Errorf("erro ao atualizar lista de compras: %w", err)
	}

	queryA := `
		SELECT lst_compras_status_name
		FROM prod.lst_compras_status
		WHERE lst_compras_status_id = $1
	`
	// Obtendo o nome do status atualizado
	err = tx.QueryRow(queryA, compras.Status_Codigo).Scan(&status_string)
	if err != nil {
		return model.LstCompras_Finish{}, err
	}

	compras.Status = status_string

	return compras, nil
}

func (r *lstComprasRepository) CancelLstCompras(compras model.LstCompras_Cancel, tx *sql.Tx) error {

	query := `
		UPDATE prod.lst_compras
		SET lst_compras_status_id = 3, lst_compras_data_atualizacao = CURRENT_TIMESTAMP
		WHERE lst_compras_id = $1
	`
	_, err := tx.Exec(query, compras.Id)
	if err != nil {
		return fmt.Errorf("erro ao cancelar lista de compras: %w", err)
	}

	return nil
}

func (r *lstComprasRepository) BuscarStatusLstCompras(lstComprasId int, tx *sql.Tx) (int, error) {
	query := `
		SELECT lst_compras_status_id 
		FROM prod.lst_compras 
		WHERE lst_compras_id = $1
	`
	var status int
	err := tx.QueryRow(query, lstComprasId).Scan(&status)
	if err != nil {
		return 0, fmt.Errorf("erro ao buscar status da lista de compras: %w", err)
	}
	return status, nil
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

	var totalItens int
	var totalPreco float64

	queryTotals, err := tx.Prepare(`
		SELECT lst_compras_total_itens, lst_compras_valor_total
		FROM prod.lst_compras
		WHERE lst_compras_id = $1
	`)

	if err != nil {
		return model.LstCompras_Post{}, fmt.Errorf("erro ao preparar a consulta para buscar totais atualizados: %v", err)
	}
	err = queryTotals.QueryRow(compras.Id).Scan(&totalItens, &totalPreco)
	if err != nil {
		return model.LstCompras_Post{}, fmt.Errorf("erro ao buscar totais atualizados: %v", err)
	}
	// Atualiza a resposta com os totais calculados

	compras.Qtd_Itens = totalItens
	compras.Total = totalPreco

	return compras, nil
}

func (r *lstComprasRepository) VerificarExistenciaLstCompras(lstComprasId int, tx *sql.Tx) (bool, error) {
	var exists bool

	query := `SELECT EXISTS(SELECT 1 FROM prod.lst_compras WHERE lst_compras_id = $1)`

	err := tx.QueryRow(query, lstComprasId).Scan(&exists)
	if err != nil {

		return false, fmt.Errorf("erro ao verificar existência da lista de compras: %v", err)
	}

	return exists, nil
}

func (r *lstComprasRepository) VerificarStatusLstCompras(lstComprasId int, tx *sql.Tx) (int, error) {
	var status int

	query := `
		SELECT lst_compras_status_id 
		FROM prod.lst_compras 
		WHERE lst_compras_id = $1
	`

	err := tx.QueryRow(query, lstComprasId).Scan(&status)
	if err != nil {

		return 0, fmt.Errorf("erro ao verificar status da lista de compras: %v", err)
	}

	return status, nil
}
