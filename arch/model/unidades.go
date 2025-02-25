package model

import "time"

type Unidades struct {
	Id               int        `json:"id" db:"unidade_id"`
	Descricao        string     `json:"descricao" db:"unidade_descricao"`
	Abreviacao       string     `json:"abreviacao" db:"unidade_abreviacao"`
	Data_cadastro    time.Time  `json:"data_cadastro" db:"unidade_data_cadastro"`
	Data_atualizacao *time.Time `json:"data_atualizacao" db:"unidade_data_atualizacao"`
}

type UnidadesPost struct {
	Id         int    `json:"id" db:"categorias_id"`
	Descricao  string `json:"descricao" db:"unidade_descricao" validate:"required,min=3,max=50"`
	Abreviacao string `json:"abreviacao" db:"unidade_abreviacao"`
}

type UnidadesUpdate struct {
	Id         int    `json:"id" db:"categorias_id"`
	Descricao  string `json:"descricao" db:"unidade_descricao" validate:"required,min=3,max=50"`
	Abreviacao string `json:"abreviacao" db:"unidade_abreviacao"`
}

type UnidadesFilters struct {
	Descricao   string `query:"descricao" db:"unidade_descricao:ILIKE:a"`	
	Limit  int    `query:"limit" validate:"min=1"`
	Page   int    `query:"page" validate:"min=1"`
}
