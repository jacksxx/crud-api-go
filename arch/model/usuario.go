package model

import (
	"errors"
	"regexp"
	"time"
)

type Usuario struct {
	Id              int        `db:"id" json:"id" validate:"required"`
	Nome            string     `db:"nome" json:"nome" validate:"required,max=250"`
	Usuario         string     `db:"usuario" json:"usuario" validate:"required,max=150"`
	Email           string     `db:"email" json:"email" validate:"required,email,max=200"`
	Senha           *string    `db:"senha" json:"senha" validate:"omitempty,max=32"`
	Salt            *string    `db:"salt" json:"salt" validate:"omitempty,max=64"`
	Ativo           bool       `db:"ativo" json:"ativo" validate:"required"`
	DataCadastro    time.Time  `db:"data_cadastro" json:"data_cadastro"`
	DataAtualizacao *time.Time `db:"data_atualizacao" json:"data_atualizacao"`
	DataInativo     *time.Time `db:"data_inativo" json:"data_inativo"`
}

type UsuarioPost struct {
	Id      int     `db:"id" json:"id"`
	Nome    string  `db:"nome" json:"nome" validate:"required,max=250"`
	Usuario string  `db:"usuario" json:"usuario" validate:"required,max=150"`
	Email   string  `db:"email" json:"email" validate:"required,email,max=200"`
	Senha   *string `db:"senha" json:"senha" validate:"required,omitempty,max=32"`
	Salt    *string `db:"salt" json:"salt" validate:"required,omitempty,max=64"`
}

type UsuarioUpdate struct {
	Id      int    `db:"id" json:"id"`
	Nome    string `db:"nome" json:"nome" validate:"required,max=250"`
	Usuario string `db:"usuario" json:"usuario" validate:"required,max=150"`
	Email   string `db:"email" json:"email" validate:"required,email,max=200"`
}

type UsuarioFilter struct {
	Ativo *bool `query:"ativo"`
	Limit int   `query:"limit" validate:"min=1"`
	Page  int   `query:"page" validate:"min=1"`
}

type SenhaUpdateRequest struct {
	SenhaAntiga string `json:"senha_antiga" validate:"required,max=32"`
	SenhaNova   string `json:"senha_nova" validate:"required,senha_complexa"`
}

func (u *UsuarioPost) Validate() []error {
	return validateNoSpaces(&u.Usuario, &u.Email)
}

func (u *UsuarioUpdate) Validate() []error {
	return validateNoSpaces(&u.Usuario, &u.Email)
}

// Função auxiliar para validação de espaços, retornando uma lista de erros
func validateNoSpaces(usuario *string, email *string) []error {
	hasSpaces := regexp.MustCompile(`\s`).MatchString
	var errorsList []error
	if usuario != nil && hasSpaces(*usuario) {
		errorsList = append(errorsList, errors.New("o campo 'Usuário' não pode ter espaços"))
	}
	if email != nil && hasSpaces(*email) {
		errorsList = append(errorsList, errors.New("o campo 'Email' não pode ter espaços"))
	}

	return errorsList
}
