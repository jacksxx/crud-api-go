package model

type Product struct {
	Id    int     `json:"id" db:"products_id" `
	Name  string  `json:"name" db:"products_name"`
	Price float64 `json:"price" db:"products_price"`
}
