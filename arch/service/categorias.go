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
	InactivateCategorias(id int) error
	ActivateCategorias(id int) error
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
	tx, err := cs.repository.BeginTransaction()
	if err != nil {
		return model.CategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = cs.repository.ValidateCategoryName(categorias.Name, &categorias.Id, tx)
	if err != nil {
		return model.CategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar categoria: %v", err)
	}

	categoria, err := cs.repository.CreateCategorias(categorias, tx)
	if err != nil {
		return model.CategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar categoria: %v", err)
	}

	// Confirmar a transação
	if err := tx.Commit(); err != nil {
		return model.CategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return categoria, http.StatusOK, nil
}

func (cs *categoriaService) UpdateCategorias(categorias model.CategoriasUpdate) (model.CategoriasUpdate, int, error) {

	tx, err := cs.repository.BeginTransaction()
	if err != nil {
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = cs.repository.ValidateCategoryName(categorias.Name, &categorias.Id, tx)
	if err != nil {
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar categoria: %v", err)
	}

	// Chama o serviço para verificar se a categoria existe
	err = cs.repository.ValidateCategory(categorias.Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar categoria: %v", err)
	}

	updatedCategories, err := cs.repository.UpdateCategoria(categorias, tx)
	if err != nil {
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar categoria: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return model.CategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return updatedCategories, http.StatusOK, nil
}

func (cs *categoriaService) InactivateCategorias(id int) error {
	tx, err := cs.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()
	// Chama o serviço para verificar se a categoria existe
	err = cs.repository.ValidateCategory(id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Inativar categoria: %v", err)
	}

	err = cs.repository.InactivateCategoria(id, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil
}

func (cs *categoriaService) ActivateCategorias(id int) error {
	tx, err := cs.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	err = cs.repository.ValidateCategory(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao Ativar categoria: %v", err)
	}

	err = cs.repository.ActivateCategoria(id, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}
