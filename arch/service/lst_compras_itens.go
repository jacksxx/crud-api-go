package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"crud-api-go/helper"
	"database/sql"
	"fmt"
	"net/http"
)

type LstComprasItensService interface {
	GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int]map[string][]model.LstCompras_Itens, error)
	GetLstComprasItensByID(id int) (model.LstCompras_Itens, int, error)
	CreateLstComprasItens(Item model.LstCompras_Itens_Post, tx *sql.Tx) (model.LstCompras_Itens_Post, int, error)
	UpdateLstComprasItem(Item model.LstCompras_Itens_Update, tx *sql.Tx) (model.LstCompras_Itens_Update, int, error)
	DeleteLstComprasItem(Item model.LstCompras_Itens_Delete, tx *sql.Tx) (int, error)
}

type lstComprasItensService struct {
	repository repository.LstComprasItensRepository
}

func NewLstComprasItensService(repository repository.LstComprasItensRepository) LstComprasItensService {
	return &lstComprasItensService{
		repository: repository,
	}
}

func (s *lstComprasItensService) GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int]map[string][]model.LstCompras_Itens, error) {
	Itens, err := s.repository.GetLstComprasItens(filters)

	if err != nil {
		return map[int]map[string][]model.LstCompras_Itens{}, fmt.Errorf("erro ao buscar produtos: %v", err)
	}

	return Itens, nil
}

func (s *lstComprasItensService) GetLstComprasItensByID(id int) (model.LstCompras_Itens, int, error) {
	Itens, err := s.repository.GetLstComprasItensById(id)
	if err != nil {
		return model.LstCompras_Itens{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar item: %v", err)
	}
	if Itens.Id == 0 {
		return model.LstCompras_Itens{}, http.StatusInternalServerError, fmt.Errorf("item não existe: %v", err)
	}
	return Itens, http.StatusOK, nil
}

func (s *lstComprasItensService) CreateLstComprasItens(Item model.LstCompras_Itens_Post, tx *sql.Tx) (model.LstCompras_Itens_Post, int, error) {

	exists, err := s.repository.CheckLstComprasExists(Item.LstCompras_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao verificar lista de compras: %v", err)
	}
	if !exists {
		return model.LstCompras_Itens_Post{}, http.StatusBadRequest, fmt.Errorf("erro: LstCompras_Id %d não encontrado", Item.LstCompras_Id)
	}
	fmt.Println("Item:", Item.LstCompras_Id)
	// Validar o produto
	err = s.repository.ValidateProduct(Item.Product_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, err
	}

	insertedItem, err := s.repository.CreateLstComprasItem(Item, tx)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar Item: %v", err)
	}

	// Atualizar os totais da lista de compras
	err = s.repository.UpdateLstComprasTotals(insertedItem.LstCompras_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar totais: %v", err)
	}

	return insertedItem, http.StatusCreated, nil
}

func (s *lstComprasItensService) UpdateLstComprasItem(Item model.LstCompras_Itens_Update, tx *sql.Tx) (model.LstCompras_Itens_Update, int, error) {
	exists, err := s.repository.CheckLstComprasExists(Item.LstCompras_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Update{}, http.StatusInternalServerError, fmt.Errorf("erro ao verificar lista de compras: %v", err)
	}
	if !exists {
		return model.LstCompras_Itens_Update{}, http.StatusBadRequest, fmt.Errorf("erro: LstCompras_Id %d não encontrado", Item.LstCompras_Id)
	}
	err = s.repository.ValidateProduct(Item.Product_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Update{}, http.StatusInternalServerError, err
	}

	updatedItem, err := s.repository.UpdateLstComprasItem(Item, tx)
	if err != nil {
		return model.LstCompras_Itens_Update{}, http.StatusInternalServerError, err
	}

	err = s.repository.UpdateLstComprasTotals(updatedItem.LstCompras_Id, tx)
	if err != nil {
		return model.LstCompras_Itens_Update{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar totais: %v", err)
	}

	return updatedItem, http.StatusOK, nil
}

func (s *lstComprasItensService) DeleteLstComprasItem(Item model.LstCompras_Itens_Delete, tx *sql.Tx) (int, error) {
	// // Verifica se o item existe antes de tentar deletar
	// var exists bool
	// queryCheck := `SELECT EXISTS(SELECT 1 FROM prod.lst_compras_itens WHERE lst_compras_itens_id = $1)`
	// err := tx.QueryRow(queryCheck, Item.Id).Scan(&exists)
	// if err != nil {
	// 	return http.StatusInternalServerError, fmt.Errorf("erro ao verificar existência do item: %v", err)
	// }
	// if !exists {
	// 	return http.StatusNotFound, fmt.Errorf("erro: item %d não encontrado na lista de compra", Item.Id)
	// }

	// Executa a remoção
	result, err := s.repository.RemoverLstComprasItem(Item.Id, tx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao executar a remoção: %v", err)
	}
	fmt.Println("id:", Item.Id, "lSTid:", Item.LstCompras_Id)
	fmt.Println("Item:", result)
	// Verifica se a remoção foi bem-sucedida
	if httpStatus, err := helper.CheckRowsAffected(result, "erro ao remover item da lista de compra"); err != nil {
		return httpStatus, err
	}

	return http.StatusOK, nil
}
