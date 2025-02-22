package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"fmt"
	"math"
	"net/http"
)

type LstComprasService interface {
	GetLstCompras(filters model.LstCompras_Filters) (model.PaginatedResponse[model.LstCompras], int, error)
	GetLstComprasById(id int) (model.LstCompras, int, error)
	CreateLstCompras(compra model.LstCompras_Post) (model.LstCompras_Post, int, error)
}

type lstComprasService struct {
	repository          repository.LstComprasRepository
	itensComprasService LstComprasItensService
}

func NewLstComprasService(repository repository.LstComprasRepository, itensComprasService LstComprasItensService) LstComprasService {
	return &lstComprasService{
		repository:          repository,
		itensComprasService: itensComprasService,
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

func (s *lstComprasService) GetLstComprasById(id int) (model.LstCompras, int, error) {
	compra, err := s.repository.GetLstComprasById(id)
	if err != nil {
		return model.LstCompras{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar lista de compra: %v", err) // Retorna erro caso a consulta falhe.
	}
	if compra.Id == 0 {
		return model.LstCompras{}, http.StatusInternalServerError, fmt.Errorf("lista de compra não existe: %v", err)
	}
	return compra, http.StatusOK, nil
}

func (s *lstComprasService) CreateLstCompras(compra model.LstCompras_Post) (model.LstCompras_Post, int, error) {
	compras, err := s.repository.CreateLstCompras(compra)
	if err != nil {
		return model.LstCompras_Post{}, http.StatusInternalServerError, fmt.Errorf("erro ao criar lista de compra: %v", err) // Retorna erro caso a inserção falhe.
	}
	// Certifique-se de que a lista de compras foi realmente criada e tem um ID válido
	if compras.Id == 0 {
		return model.LstCompras_Post{}, http.StatusInternalServerError, fmt.Errorf("erro: ID da lista de compras inválido")
	}
	var itensCriados []model.LstCompras_Itens_Post
	for _, material := range compras.LstCompras_Itens {
		material.LstCompras_Id = compras.Id
		ItensCompras, httpStatus, err := s.itensComprasService.CreateLstComprasItens(material)

		if err != nil {
			return model.LstCompras_Post{}, http.StatusInternalServerError, err
		}

		if httpStatus != http.StatusCreated {
			return model.LstCompras_Post{}, http.StatusInternalServerError, err

		}

		err = s.itensComprasService.ValidateProduct(ItensCompras.Product_Id)
		if err != nil {
			return model.LstCompras_Post{}, http.StatusInternalServerError, err
		}

		itensCriados = append(itensCriados, ItensCompras)
	}

	compras.LstCompras_Itens = itensCriados

	return compras, http.StatusCreated, nil

}
