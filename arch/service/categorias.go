package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
	"net/http"
)

type CategoriaService interface {
	GetCategorias(filters model.CategoriasFilters) (model.PaginatedResponse[model.Categorias], int, error)
	GetProductByID(id int) (model.Categorias, int, error)
	CreateCategorias(categorias model.CategoriasPost) (model.CategoriasPost, int, error)
	UpdateCategorias(categorias model.CategoriasUpdate) (model.CategoriasUpdate, int, error)
	DeleteCategorias(id int) error
	ValidateCategoryName(nomeCategoria string, categoriaId *int) error
	ValidateCategory(categoriaId int) error
}

type categoriaService struct {
	repository repository.CategoriasRepository
}

func NewCategoriaService(repository repository.CategoriasRepository) CategoriaService {
	return &categoriaService{
		repository: repository,
	}
}

func (cs *categoriaService) GetCategorias(filters model.CategoriasFilters) (model.PaginatedResponse[model.Categorias], int, error) {
	categories, err := cs.repository.GetCategorias(filters)
	if err != nil {
		return model.PaginatedResponse[model.Categorias]{}, 0, fmt.Errorf("erro ao buscar categorias: %v", err)
	}
	total, err := cs.repository.CountCategories(filters)
	if err != nil {
		return model.PaginatedResponse[model.Categorias]{}, 0, fmt.Errorf("erro ao buscar quantidade de categorias: %v", err)
	}
	response := model.PaginatedResponse[model.Categorias]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       categories,
	}
	return response, total, nil
}

func (cs *categoriaService) GetProductByID(id int) (model.Categorias, int, error) {
	categorias, err := cs.repository.GetCategoriasById(id)
	if err != nil {
		return model.Categorias{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar categoria: %v", err)
	}
	if categorias.Id == 0 {
		return model.Categorias{}, http.StatusInternalServerError, fmt.Errorf("categoria não existe: %v", err)
	}
	return categorias, http.StatusOK, nil
}

func (cs *categoriaService) CreateCategorias(categorias model.CategoriasPost) (model.CategoriasPost, int, error) {

	categoriaID, err := cs.repository.CreateCategorias(categorias)
	if err != nil {
		return model.CategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar categoria: %v", err)
	}

	categorias.Id = categoriaID

	return categorias, http.StatusOK, nil
}

func (cs *categoriaService) UpdateCategorias(categorias model.CategoriasUpdate) (model.CategoriasUpdate, int, error) {

	updatedCategories, err := cs.repository.UpdateCategoria(categorias)
	if err != nil {
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar categoria: %v", err)
	}
	return updatedCategories, http.StatusOK, nil
}

func (cs *categoriaService) DeleteCategorias(id int) error {
	err := cs.repository.DeleteCategoria(id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *categoriaService) ValidateCategoryName(nomeCategoria string, categoriaId *int) error {
	err := cs.repository.ValidateCategoryName(nomeCategoria, categoriaId)
	if err != nil {
		return err
	}
	return nil
}

func (cs *categoriaService) ValidateCategory(categoriaId int) error {
	err := cs.repository.ValidateCategory(categoriaId)
	if err != nil {
		return err
	}
	return nil
}
