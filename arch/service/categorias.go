package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
)

type CategoriaService interface {
	GetCategorias(filters model.CategoriasFilters) ([]model.Categorias, error)
	GetProductByID(id int) (*model.Categorias, error)
	CreateCategorias(categorias model.CategoriasPost) (model.CategoriasPost, error)
	UpdateCategorias(categorias model.CategoriasUpdate) (model.CategoriasUpdate, error)
	DeleteCategorias(id int) error
	ValidateCategoryName(nomeCategoria string) error
}

type categoriaService struct {
	repository repository.CategoriasRepository
}

func NewCategoriaService(repository repository.CategoriasRepository) CategoriaService {
	return &categoriaService{
		repository: repository,
	}
}

func (cs *categoriaService) GetCategorias(filters model.CategoriasFilters) ([]model.Categorias, error) {
	return cs.repository.GetCategorias(filters)
}

func (cs *categoriaService) GetProductByID(id int) (*model.Categorias, error) {
	categorias, err := cs.repository.GetCategoriasById(id)
	if err != nil {
		return nil, err
	}

	return categorias, nil
}

func (cs *categoriaService) CreateCategorias(categorias model.CategoriasPost) (model.CategoriasPost, error) {

	categoriaID, err := cs.repository.CreateCategorias(categorias)
	if err != nil {
		return model.CategoriasPost{}, err
	}

	categorias.Id = categoriaID

	return categorias, nil
}

func (cs *categoriaService) UpdateCategorias(categorias model.CategoriasUpdate) (model.CategoriasUpdate, error) {

	updatedCategories, err := cs.repository.UpdateCategoria(categorias)
	if err != nil {
		return model.CategoriasUpdate{}, err
	}
	return updatedCategories, nil
}

func (cs *categoriaService) DeleteCategorias(id int) error {
	err := cs.repository.DeleteCategoria(id)
	if err != nil {
		return err
	}
	return nil
}

func (cs *categoriaService) ValidateCategoryName(nomeCategoria string) error {
	err := cs.repository.ValidateCategoryName(nomeCategoria)
	if err != nil {
		return err
	}
	return nil
}
