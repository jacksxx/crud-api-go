package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
)

type SubcategoriasService interface {
	GetSubcategorias(filters model.SubcategoriasFilters) (model.PaginatedResponse[model.Subcategorias], int, error)
}

type subcategoriaService struct {
	repository repository.SubcategoriasRepository
}

func NewSubcategoriaService(repository repository.SubcategoriasRepository) SubcategoriasService {
	return &subcategoriaService{
		repository: repository,
	}
}

func (s *subcategoriaService) GetSubcategorias(filters model.SubcategoriasFilters) (model.PaginatedResponse[model.Subcategorias], int, error) {
	subcategorias, err := s.repository.GetSubcategorias(filters)
	if err != nil {
		return model.PaginatedResponse[model.Subcategorias]{}, 0, fmt.Errorf("erro ao buscar categorias: %v", err)
	}
	total, err := s.repository.CountSubcategories(filters)
	if err != nil {
		return model.PaginatedResponse[model.Subcategorias]{}, 0, fmt.Errorf("erro ao buscar quantidade de categorias: %v", err)
	}
	response := model.PaginatedResponse[model.Subcategorias]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       subcategorias,
	}
	return response, total, nil
}
