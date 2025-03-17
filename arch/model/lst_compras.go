package model

import "time"

type LstCompras struct {
	Id               int        `json:"id" db:"lst_compras_id"`
	Nome             string     `json:"nome" db:"lst_compras_name"`
	Total            float64    `json:"total" db:"lst_compras_valor_total"`
	Qtd_Itens        int        `json:"qtd_itens" db:"lst_compras_total_itens"`
	Status_Codigo    int        `json:"status_codigo" db:"lst_compras_status_id"`
	Status           string     `json:"status" db:"lst_compras_status_name"`
	Data_Cadastro    time.Time  `json:"data_cadastro" db:"lst_compras_data_cadastro"`
	Data_Atualizacao *time.Time `json:"data_atualizacao" db:"lst_compras_data_atualizacao"`
}

type LstCompras_Post struct {
	Id               int                     `json:"id" db:"lst_compras_id"`
	Nome             string                  `json:"nome" db:"lst_compras_name"`
	Status_Codigo    int                     `json:"status_codigo" db:"lst_compras_status_id"`
	Status           string                  `json:"status" db:"lst_compras_status_name"`
	Total            float64                 `json:"total" db:"lst_compras_valor_total"`
	Qtd_Itens        int                     `json:"qtd_itens" db:"lst_compras_total_itens"`
	LstCompras_Itens []LstCompras_Itens_Post `json:"lstcompras_itens" validate:"required,dive"`
}

type LstCompras_Update struct {
	Id               int                       `json:"id" db:"lst_compras_id"`
	Nome             string                    `json:"nome" db:"lst_compras_name"`
	Status_Codigo    int                       `json:"status_codigo" db:"lst_compras_status_id"`
	Status           string                    `json:"status" db:"lst_compras_status_name"`
	Total            float64                   `json:"total" db:"lst_compras_valor_total"`
	Qtd_Itens        int                       `json:"qtd_itens" db:"lst_compras_total_itens"`
	LstCompras_Itens []LstCompras_Itens_Update `json:"lstcompras_itens" validate:"required,dive"`
}

type LstCompras_Finish struct {
	Id                             int                        `json:"id" db:"lst_compras_id"`
	Nome                           string                     `json:"nome" db:"lst_compras_name"`
	Status_Codigo                  int                        `json:"status_codigo" db:"lst_compras_status_id"`
	Status                         string                     `json:"status" db:"lst_compras_status_name"`
	Total                          float64                    `json:"total" db:"lst_compras_valor_total"`
	Qtd_Itens                      int                        `json:"qtd_itens" db:"lst_compras_total_itens"`
	LstCompras_Itens_Comprados     []LstCompras_Itens_Finish  `json:"lstcompras_itens_comprados" validate:"required,dive"`
	LstCompras_Itens_Nao_Comprados *[]LstCompras_Itens_Finish `json:"lstcompras_itens_nao_comprados" validate:"required,dive"`
}

type LstCompras_Filters struct {
	Nome          string `query:"nome" db:"lst_compras_name:ILIKE:a"`
	Status_Codigo int    `query:"status_codigo" db:"lst_compras_status_id"`
	Limit         int    `query:"limit" validate:"min=1"`
	Page          int    `query:"page" validate:"min=1"`
}
