package model

import "time"

type LstCompras_Itens struct {
	Id                 int        `json:"id" db:"lst_compras_itens_id"`
	LstCompras_Id      int        `json:"lst_compras_id" db:"lst_compras_id"`
	Product_Id         int        `json:"product_id" db:"products_id"`
	Product_Name       string     `json:"product_name" db:"products_name"`
	Unidade_Descricao  string     `json:"unidade_descricao" db:"unidade_descricao"`
	Unidade_Abreviacao string     `json:"unidade_abreviacao" db:"unidade_abreviacao"`
	Quantidade         int        `json:"quantidade" db:"lst_compras_itens_quantidade"`
	Preco              float64    `json:"preco" db:"lst_compras_itens_preco"`
	Item_Check         bool       `json:"item_check" db:"lst_compras_itens_comprado"`
	Data_Cadastro      time.Time  `json:"data_cadastro" db:"lst_compras_itens_data_cadastro"`
	Data_Atualizacao   *time.Time `json:"data_atualizacao" db:"lst_compras_itens_data_atualizacao"`
}

type LstCompras_Itens_Post struct {
	Id                 int     `json:"id" db:"lst_compras_itens_id"`
	LstCompras_Id      int     `json:"lst_compras_id" db:"lst_compras_id"`
	Product_Name       string  `json:"product_name" db:"products_name"`
	Product_Id         int     `json:"product_id" db:"products_id" validate:"required"`
	Unidade_Descricao  string  `json:"unidade_descricao" db:"unidade_descricao"`
	Unidade_Abreviacao string  `json:"unidade_abreviacao" db:"unidade_abreviacao"`
	Quantidade         int     `json:"quantidade" db:"lst_compras_itens_quantidade" validate:"required"`
	Preco              float64 `json:"preco" db:"lst_compras_itens_preco" validate:"required"`
}

type LstCompras_Itens_Update struct {
	Id            int     `json:"id" db:"lst_compras_itens_id"`
	LstCompras_Id int     `json:"lst_compras_id" db:"lst_compras_id"`
	Product_Id    int     `json:"product_id" db:"products_id" validate:"required"`
	Quantidade    int     `json:"quantidade" db:"lst_compras_itens_quantidade" validate:"required"`
	Preco         float64 `json:"preco" db:"lst_compras_itens_preco" validate:"required"`
	Acao          string  `json:"acao" validate:"required,oneof=adicionar remover atualizar"`
}

type LstCompras_Itens_Finish struct {
	Id                 int     `json:"id" db:"lst_compras_itens_id"`
	LstCompras_Id      int     `json:"lst_compras_id" db:"lst_compras_id"`
	Product_Name       string  `json:"product_name" db:"products_name"`
	Product_Id         int     `json:"product_id" db:"products_id"`
	Unidade_Descricao  string  `json:"unidade_descricao" db:"unidade_descricao"`
	Unidade_Abreviacao string  `json:"unidade_abreviacao" db:"unidade_abreviacao"`
	Quantidade         int     `json:"quantidade" db:"lst_compras_itens_quantidade" validate:"required"`
	Preco              float64 `json:"preco" db:"lst_compras_itens_preco" validate:"required"`
	Item_Check         *bool   `json:"item_check" db:"lst_compras_itens_comprado" validate:"required"`
}

type LstCompras_Itens_Delete struct {
	Id            int `json:"id" db:"lst_compras_itens_id"`
	LstCompras_Id int `json:"lst_compras_id" db:"lst_compras_id"`
}

type LstCompras_Itens_Filters struct {
	LstCompras_Id int    `query:"lst_compras_id"`
	Product_Id    int    `query:"products_id"`
	Product_Name  string `query:"products_name" db:"products_name:ILIKE:a"`
	Limit         int    `query:"limit" validate:"min=1"`
	Page          int    `query:"page" validate:"min=1"`
}
