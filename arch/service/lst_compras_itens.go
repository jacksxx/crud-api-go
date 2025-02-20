package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
)

type LstComprasItensService interface {
	GetLstComprasItens(filters model.LstCompras_Itens_Filters) ([][]model.LstCompras_Itens, error)
}

type lstComprasItensService struct {
	repository repository.LstComprasItensRepository
}

func NewLstComprasItensService(repository repository.LstComprasItensRepository) LstComprasItensService {
	return &lstComprasItensService{
		repository: repository,
	}
}

func (s *lstComprasItensService) GetLstComprasItens(filters model.LstCompras_Itens_Filters) ([][]model.LstCompras_Itens, error) {
	Itens, err := s.repository.GetLstComprasItens(filters)

	if err != nil {
		return [][]model.LstCompras_Itens{}, fmt.Errorf("erro ao buscar produtos: %v", err)
	}

	return Itens, nil
}
