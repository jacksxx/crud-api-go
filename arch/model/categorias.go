package model

import "time"

type Categorias struct {
	Id               int        `json:"id" db:"categorias_id"`
	Name             string     `json:"nome" db:"categorias_name"`
	Data_cadastro    time.Time  `json:"data_cadastro" db:"categorias_data_cadastro"`
	Data_atualizacao *time.Time `json:"data_atualizacao" db:"categorias_data_atualizacao"`
	Data_Inativacao  *time.Time `json:"data_inativacao" db:"categorias_data_inativacao"`
	Status           string     `json:"status" db:"categorias_status"`
}

type CategoriasPost struct {
	Id   int    `json:"id" db:"categorias_id"`
	Name string `json:"nome" db:"categorias_name" validate:"required,min=3,max=50"`
}

type CategoriasUpdate struct {
	Id   int    `json:"id" db:"categorias_id"`
	Name string `json:"nome" db:"categorias_name" validate:"required,min=3,max=50"`
}

type CategoriasFilters struct {
	Name   string `query:"nome" db:"categorias_name:ILIKE:a"`
	Status string `query:"status"`
	Limit  int    `query:"limit" validate:"min=1"`
	Page   int    `query:"page" validate:"min=1"`
}
