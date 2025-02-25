package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
	"net/http"
)

type UnidadesService interface {
	GetUnits(filters model.UnidadesFilters) (model.PaginatedResponse[model.Unidades], int, error)
	GetUnitByID(id int) (model.Unidades, int, error)
	CreateUnit(unit model.UnidadesPost) (model.UnidadesPost, int, error)
	UpdateUnit(unit model.UnidadesUpdate) (model.UnidadesUpdate, int, error)
}

type unidadesService struct {
	repository repository.UnidadesRepository
}

func NewUnidadesService(repository repository.UnidadesRepository) UnidadesService {
	return &unidadesService{
		repository: repository,
	}
}

func (s *unidadesService) GetUnits(filters model.UnidadesFilters) (model.PaginatedResponse[model.Unidades], int, error) {
	unidades, err := s.repository.GetUnidades(filters)
	if err != nil {
		return model.PaginatedResponse[model.Unidades]{}, 0, fmt.Errorf("erro ao buscar unidades: %v", err)
	}
	total, err := s.repository.CountUnits(filters)
	if err != nil {
		return model.PaginatedResponse[model.Unidades]{}, 0, fmt.Errorf("erro ao buscar quantidade de unidades: %v", err)
	}
	response := model.PaginatedResponse[model.Unidades]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       unidades,
	}
	return response, total, nil
}

func (s *unidadesService) GetUnitByID(id int) (model.Unidades, int, error) {
	unit, err := s.repository.GetUnidadesById(id)
	if err != nil {
		return model.Unidades{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar unidade: %v", err)
	}
	if unit.Id == 0 {
		return model.Unidades{}, http.StatusInternalServerError, fmt.Errorf("unidade não existe: %v", err)
	}
	return unit, http.StatusOK, nil
}

func (s *unidadesService) CreateUnit(unit model.UnidadesPost) (model.UnidadesPost, int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.UnidadesPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	units, err := s.repository.CreateUnidades(unit, tx)
	if err != nil {
		return model.UnidadesPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar unidade: %v", err)
	}
	err = s.repository.ValidateUnitName(units.Descricao, &units.Id)
	if err != nil {
		return model.UnidadesPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar unidade: %v", err.Error())
	}
	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return model.UnidadesPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return units, http.StatusOK, nil
}

func (s *unidadesService) UpdateUnit(unit model.UnidadesUpdate) (model.UnidadesUpdate, int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.UnidadesUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()
	
	units, err := s.repository.UpdateUnidades(unit, tx)
	if err != nil {
		return model.UnidadesUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao alterar unidade: %v", err)
	}

	err = s.repository.ValidateUnit(units.Id)
	if err != nil {
		return model.UnidadesUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao alterar unidade: %v", err.Error())
	}

	err = s.repository.ValidateUnitName(units.Descricao, &units.Id)
	if err != nil {
		return model.UnidadesUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao alterar unidade: %v", err.Error())
	}

	// Confirmar a transação
	err = tx.Commit()
	if err != nil {
		return model.UnidadesUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return units, http.StatusOK, nil
}
