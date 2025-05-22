package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
	"crud-api-go/helper"
	"fmt"
	"math"
	"net/http"
)

type UsuarioService interface {
	GetUsuarios(filters model.UsuarioFilter) (model.PaginatedResponse[model.Usuario], error)
	GetUsuarioById(id int) (model.Usuario, int, error)
	CreateUsuario(usuario model.UsuarioPost) (model.UsuarioPost, int, error)
	UpdateUsuario(usuario model.UsuarioUpdate) (model.UsuarioUpdate, int, error)
	InativarUsuario(id int) error
	AtivarUsuario(id int) error
	UpdateSenha(id int, dados model.SenhaUpdateRequest) (int, error)
	ResetarSenha(id int) (int, error)
}

type usuarioService struct {
	repository repository.UsuarioRepository
}

func NewUsuarioService(repository repository.UsuarioRepository) UsuarioService {
	return &usuarioService{
		repository: repository,
	}
}

func (s *usuarioService) GetUsuarios(filters model.UsuarioFilter) (model.PaginatedResponse[model.Usuario], error) {
	usuarios, err := s.repository.GetUsuarios(filters)
	if err != nil {
		return model.PaginatedResponse[model.Usuario]{}, err
	}
	total, err := s.repository.CountUsuario(filters)
	if err != nil {
		return model.PaginatedResponse[model.Usuario]{}, fmt.Errorf("erro ao buscar quantidade de categorias: %v", err)
	}
	response := model.PaginatedResponse[model.Usuario]{
		Total:      total,
		Page:       filters.Page,
		TotalPages: int(math.Ceil(float64(total) / float64(filters.Limit))),
		Data:       usuarios,
	}
	return response, nil
}

func (s *usuarioService) GetUsuarioById(id int) (model.Usuario, int, error) {
	usuario, err := s.repository.GetUsuariosById(id)
	if err != nil {
		return model.Usuario{}, http.StatusInternalServerError, fmt.Errorf("erro ao buscar usuário: %v", err)
	}
	if usuario.Id == 0 {
		return model.Usuario{}, http.StatusNotFound, fmt.Errorf("usuário não encontrado com o código: %d", id)
	}
	return usuario, http.StatusOK, nil
}

func (s *usuarioService) CreateUsuario(usuario model.UsuarioPost) (model.UsuarioPost, int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	exists, err := s.repository.CheckEmail(usuario.Email, tx, nil)
	if err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, err
	}
	if exists {
		return model.UsuarioPost{}, http.StatusConflict, fmt.Errorf("já existe um usuário com esse email")
	}

	exists, err = s.repository.CheckUsuario(usuario.Usuario, tx, nil)
	if err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, err
	}
	if exists {
		return model.UsuarioPost{}, http.StatusConflict, fmt.Errorf("já existe um usuário com esse nome de usuário")
	}

	hash, salt, err := helper.HashPassword(*usuario.Senha)
	if err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, err
	}

	user, err := s.repository.CreateUsuario(tx, usuario, hash, salt)
	if err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, err
	}

	if err := tx.Commit(); err != nil {
		return model.UsuarioPost{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return user, http.StatusCreated, nil
}

func (s *usuarioService) UpdateUsuario(usuario model.UsuarioUpdate) (model.UsuarioUpdate, int, error) {
	if usuario.Nome == "" && usuario.Email == "" && usuario.Usuario == "" {
		return model.UsuarioUpdate{}, http.StatusBadRequest, fmt.Errorf("nenhum campo fornecido para atualização")
	}
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()
	isActive, err := s.repository.CheckUsuarioAtivo(usuario.Id, tx)
	if err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, err
	}
	if !isActive {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, fmt.Errorf("o usuário não está ativo")
	}
	exists, err := s.repository.CheckEmail(usuario.Email, tx, &usuario.Id)
	if err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, err
	}
	if exists {
		return model.UsuarioUpdate{}, http.StatusConflict, fmt.Errorf("já existe um usuário com esse email")
	}

	exists, err = s.repository.CheckUsuario(usuario.Usuario, tx, &usuario.Id)
	if err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, err
	}
	if exists {
		return model.UsuarioUpdate{}, http.StatusConflict, fmt.Errorf("já existe um usuário com esse nome de usuário")
	}

	user, err := s.repository.UpdateUsuario(tx, usuario)
	if err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao executar a atualização: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return model.UsuarioUpdate{}, http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}

	return user, http.StatusOK, nil
}

func (s *usuarioService) InativarUsuario(id int) error {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	ativo, err := s.repository.CheckUsuarioAtivo(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao verificar usuário")
	}
	if !ativo {
		return fmt.Errorf("usuário já está inativo")
	}

	err = s.repository.InativarUsuario(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao executar a inativação: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil
}

func (s *usuarioService) AtivarUsuario(id int) error {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	ativo, err := s.repository.CheckUsuarioAtivo(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao verificar usuário")
	}
	if ativo {
		return fmt.Errorf("usuário já está ativo")
	}

	err = s.repository.AtivarUsuario(id, tx)
	if err != nil {
		return fmt.Errorf("erro ao executar a ativação: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return nil
}

func (s *usuarioService) UpdateSenha(id int, dados model.SenhaUpdateRequest) (int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	hashRecord, saltRecord, err := s.repository.GetSenhaNSalt(id, tx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao buscar usuário: %v", err)
	}
	if hashRecord == "" || saltRecord == "" {
		return http.StatusNotFound, fmt.Errorf("usuário não encontrado")
	}

	valid := helper.CheckPassword(hashRecord, saltRecord, dados.SenhaAntiga)
	if !valid {
		return http.StatusBadRequest, fmt.Errorf("senha antiga inválida")
	}

	hash, salt, err := helper.HashPassword(dados.SenhaNova)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao atualizar senha: %v", err)
	}

	err = s.repository.UpdateSenha(id, hash, salt, tx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao atualizar senha: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return http.StatusOK, nil
}

func (s *usuarioService) ResetarSenha(id int) (int, error) {
	tx, err := s.repository.BeginTransaction()
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao iniciar transação: %v", err)
	}
	defer tx.Rollback()

	usuario, err := s.repository.GetUsuariosById(id)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao buscar usuário")
	}
	if usuario.Id == 0 {
		return http.StatusNotFound, fmt.Errorf("usuário não encontrado ")
	}

	novaSenha := helper.GerarSenhaPadrao(usuario.Nome)

	hash, salt, err := helper.HashPassword(novaSenha)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao resetar senha: %v", err)
	}

	err = s.repository.UpdateSenha(id, hash, salt, tx)
	if err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao resetar senha: %v", err)
	}

	if err := tx.Commit(); err != nil {
		return http.StatusInternalServerError, fmt.Errorf("erro ao confirmar transação: %v", err)
	}
	return http.StatusOK, nil
}
