package model

import "time"

type Subcategorias struct {
	Id               int        `json:"id" db:"subcategorias_id"`
	Name             string     `json:"nome" db:"subcategorias_name"`
	CategoriasId     int        `json:"categorias_id" db:"categorias_id"`
	CategoriasName   string     `json:"categorias_nome" db:"categorias_name"`
	Data_cadastro    time.Time  `json:"data_cadastro" db:"subcategorias_data_cadastro"`
	Data_atualizacao *time.Time `json:"data_atualizacao" db:"subcategorias_data_atualizacao"`
	Data_Inativacao  *time.Time `json:"data_inativacao" db:"subcategorias_data_inativacao"`
	Status           string     `json:"status" db:"subcategorias_status"`
}

type SubcategoriasPost struct {
	Id               int        `json:"id" db:"subcategorias_id"`
	Name             string     `json:"nome" db:"subcategorias_name"`
	CategoriasId     int        `json:"categorias_id" db:"categorias_id"`	
}

type SubcategoriasUpdate struct {
	Id               int        `json:"id" db:"subcategorias_id"`
	Name             string     `json:"nome" db:"subcategorias_name"`
	CategoriasId     int        `json:"categorias_id" db:"categorias_id"`	
}

type SubcategoriasFilters struct {
	Name   string `query:"nome" db:"subcategorias_name:ILIKE:a"`
	Status string `query:"status"`
	Limit  int    `query:"limit" validate:"min=1"`
	Page   int    `query:"page" validate:"min=1"`
}
