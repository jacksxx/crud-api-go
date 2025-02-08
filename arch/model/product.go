package model

import "time"

type Product struct {
	Id             int       `json:"id" db:"products_id" `
	Name           string    `json:"name" db:"products_name"`
	Price          float64   `json:"price" db:"products_price"`
	Categoria_Id   int       `json:"categorias_id" db:"categorias_id"`
	Categoria_Name string    `json:"categoria_name" db:"categorias_name"`
	Data_Cadastro  time.Time `json:"data_cadastro" db:"data_cadastro"`
}

type ProductFilters struct {
	Name         string `query:"name" db:"products_name:ILIKE:a"`
	Categoria_Id int    `query:"categorias_id"`
	Limit        int    `query:"limit" validate:"min=1"`
	Page         int    `query:"page" validate:"min=1"`
}
