package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
)

type LstComprasService interface {
	GetLstCompras(filters model.LstCompras_Filters) (model.PaginatedResponse[model.LstCompras], int, error)
}

type lstComprasService struct {
	repository repository.LstComprasRepository
}

func NewLstComprasService(repository repository.LstComprasRepository) LstComprasService {
	return &lstComprasService{
		repository: repository,
	}
}

func (s *lstComprasService) GetLstCompras(filters model.LstCompras_Filters) (model.PaginatedResponse[model.LstCompras], int, error) {
	Compras, err := s.repository.GetLstCompras(filters)
	if err != nil {
		return model.PaginatedResponse[model.LstCompras]{}, 0, fmt.Errorf("erro ao buscar lista de compras: %v", err)
	}

	total, err := s.repository.CountLstCompras(filters)
	if err != nil {
		return model.PaginatedResponse[model.LstCompras]{}, 0, fmt.Errorf("erro ao buscar quantidade de lista de compras: %v", err)
	}

	response := model.PaginatedResponse[model.LstCompras]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       Compras,
	}

	return response, total, nil
}
