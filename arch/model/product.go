package model

import "time"

type Product struct {
	Id                int        `json:"id" db:"products_id"`
	Name              string     `json:"name" db:"products_name"`
	Price             float64    `json:"price" db:"products_price"`
	Categoria_Id      int        `json:"categorias_id" db:"categorias_id"`
	Categoria_Name    string     `json:"categoria_name" db:"categorias_name"`
	Unidade_Id        int        `json:"unidade_id" db:"unidade_id"`
	Unidade_Descricao string     `json:"unidade_descricao" db:"unidade_descricao"`
	Unidade_Abreviacao string     `json:"unidade_abreviacao" db:"unidade_abreviacao"`
	Data_Cadastro     time.Time  `json:"data_cadastro" db:"products_data_cadastro"`
	Data_Atualizacao  *time.Time `json:"data_atualizacao" db:"products_data_atualizacao"`
	Data_Inativacao   *time.Time `json:"data_inativacao" db:"products_data_inativacao"`
	Status            string     `json:"status" db:"products_status"`
}

type ProductPost struct {
	Id                int     `json:"id" db:"products_id"` // Pode ser omitido se for gerado pelo banco
	Name              string  `json:"name" db:"products_name" validate:"required,min=3,max=100"`
	Price             float64 `json:"price" db:"products_price" validate:"required,gt=0"`
	Categoria_Id      int     `json:"categorias_id" db:"categorias_id" validate:"required,gt=0"`
	Categoria_Name    string  `json:"categoria_name" db:"categorias_name"`
	Unidade_Id        int     `json:"unidade_id" db:"unidade_id"`
	Unidade_Descricao string  `json:"unidade_descricao" db:"unidade_descricao"`
}

type ProductUpdate struct {
	Id                int     `json:"id" db:"products_id"`
	Name              string  `json:"name" db:"products_name" validate:"omitempty,min=3,max=100"`
	Price             float64 `json:"price" db:"products_price" validate:"omitempty,gt=0"`
	Categoria_Id      int     `json:"categorias_id" db:"categorias_id" validate:"omitempty,gt=0"`
	Categoria_Name    string  `json:"categoria_name" db:"categorias_name"`
	Unidade_Id        int     `json:"unidade_id" db:"unidade_id"`
	Unidade_Descricao string  `json:"unidade_descricao" db:"unidade_descricao"`
}

type ProductFilters struct {
	Name         string `query:"name" db:"products_name:ILIKE:a"`
	Categoria_Id int    `query:"categorias_id"`
	Unidade_Id   int    `query:"unidade_id"`
	Status       string `query:"status"`
	Limit        int    `query:"limit" validate:"min=1"`
	Page         int    `query:"page" validate:"min=1"`
}
