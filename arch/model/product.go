package model

type Product struct {
	Id    int     `json:"id" db:"products_id" `
	Name  string  `json:"name" db:"products_name"`
	Price float64 `json:"price" db:"products_price"`
}

type ProductFilters struct {
	Name  string `query:"name" db:"products_name:ILIKE:a"`
	Limit int    `query:"limit" validate:"min=1"`
	Page  int    `query:"page" validate:"min=1"`
}
