package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
	"net/http"
)

type SubcategoriasService interface {
	GetSubcategorias(filters model.SubcategoriasFilters) (model.PaginatedResponse[model.Subcategorias], int, error)
	GetSubcategoriasById(id int) (model.Subcategorias, int, error)
	CreateSubcategorias(subcategorias model.SubcategoriasPost) (model.SubcategoriasPost, int, error)
	UpdateSubcategorias(subcategorias model.SubcategoriasUpdate) (model.SubcategoriasUpdate, int, error)
	InactivateSubategorias(id int) error
	ActivateSubategorias(id int) error
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

func (s *subcategoriaService) GetSubcategoriasById(id int) (model.Subcategorias, int, error) {
	subcategorias, err := s.repository.GetSubcategoriasById(id)
	if err != nil {
		return model.Subcategorias{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar categorias: %v", err)
	}
	if subcategorias.Id == 0 {
		return model.Subcategorias{}, http.StatusInternalServerError, fmt.Errorf("categoria não existe: %v", err)
	}

	return subcategorias, http.StatusOK, nil
}

func (s *subcategoriaService) CreateSubcategorias(subcategorias model.SubcategoriasPost) (model.SubcategoriasPost, int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.SubcategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	err = s.repository.ValidateSubCategoryName(subcategorias.Name, &subcategorias.Id, tx)
	if err != nil {
		return model.SubcategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar subcategoria: %v", err)
	}

	err = s.repository.ValidateCategory(subcategorias.CategoriasId, tx)
	if err != nil {
		return model.SubcategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar subcategoria: %v", err)
	}

	subcategoria, err := s.repository.CreateSubcategorias(subcategorias, tx)
	if err != nil {
		return model.SubcategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar subcategoria: %v", err)
	}

	// Confirmar a transação
	if err := tx.Commit(); err != nil {
		return model.SubcategoriasPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return subcategoria, http.StatusOK, nil
}

func (s *subcategoriaService) UpdateSubcategorias(subcategorias model.SubcategoriasUpdate) (model.SubcategoriasUpdate, int, error) {

	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	// Chama o serviço para verificar se a categoria existe
	err = s.repository.ValidateSubcategory(subcategorias.Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar subcategoria: %v", err)
	}

	err = s.repository.ValidateSubCategoryName(subcategorias.Name, &subcategorias.Id, tx)
	if err != nil {
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar subcategoria: %v", err)
	}

	err = s.repository.ValidateCategory(subcategorias.CategoriasId, tx)
	if err != nil {
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar subcategoria: %v", err)
	}

	status, err := s.repository.CheckStatus(subcategorias.Id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar subcategoria: %v", err)
	}
	if status == "inativo" {
		// Retorna erro 400 caso o ID da categoria não exista
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar subcategoria: Subcategoria já inativada")
	}

	updatedSubcategories, err := s.repository.UpdateSubcategorias(subcategorias, tx)
	if err != nil {
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao atualizar subcategoria: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return model.SubcategoriasUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return updatedSubcategories, http.StatusOK, nil
}

func (s *subcategoriaService) InactivateSubategorias(id int) error {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()
	// Chama o serviço para verificar se a categoria existe
	err = s.repository.ValidateSubcategory(id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Inativar subcategoria: %v", err)
	}

	status, err := s.repository.CheckStatus(id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Inativar subcategoria: %v", err)
	}
	if status == "inativo" {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Inativar subcategoria: Subcategoria já inativada")
	}

	err = s.repository.InactivateSubCategoria(id, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil
}

func (s *subcategoriaService) ActivateSubategorias(id int) error {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}

	err = s.repository.ValidateSubcategory(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao Ativar subcategoria: %v", err)
	}

	status, err := s.repository.CheckStatus(id, tx)
	if err != nil {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Ativar subcategoria: %v", err)
	}
	if status == "ativo" {
		// Retorna erro 400 caso o ID da categoria não exista
		return fmt.Errorf("erro ao Ativar subcategoria: Subcategoria já ativada")
	}

	err = s.repository.ActivateSubCategoria(id, tx)
	if err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return nil
}
