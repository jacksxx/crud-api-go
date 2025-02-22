package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"net/http"
)

type LstComprasItensService interface {
	GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int][]model.LstCompras_Itens, error)
	GetLstComprasItensByID(id int) (model.LstCompras_Itens, int, error)
	CreateLstComprasItens(Item model.LstCompras_Itens_Post) (model.LstCompras_Itens_Post, int, error)
	ValidateProduct(productId int) error
}

type lstComprasItensService struct {
	repository repository.LstComprasItensRepository
}

func NewLstComprasItensService(repository repository.LstComprasItensRepository) LstComprasItensService {
	return &lstComprasItensService{
		repository: repository,
	}
}

func (s *lstComprasItensService) GetLstComprasItens(filters model.LstCompras_Itens_Filters) (map[int][]model.LstCompras_Itens, error) {
	Itens, err := s.repository.GetLstComprasItens(filters)

	if err != nil {
		return map[int][]model.LstCompras_Itens{}, fmt.Errorf("erro ao buscar produtos: %v", err)
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

func (s *lstComprasItensService) CreateLstComprasItens(Item model.LstCompras_Itens_Post) (model.LstCompras_Itens_Post, int, error) {
	itemId, err := s.repository.CreateLstComprasItem(Item)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar Item: %v", err)
	}

	Item.Id = itemId

	// Atualizar os totais da lista de compras
	err = s.repository.UpdateLstComprasTotals(Item.LstCompras_Id)
	if err != nil {
		return model.LstCompras_Itens_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar totais: %v", err)
	}

	return Item, http.StatusCreated, nil
}

func (s *lstComprasItensService) ValidateProduct(productId int) error {
	err := s.repository.ValidateProduct(productId)
	if err != nil {
		return err
	}
	return nil
}
