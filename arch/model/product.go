package model

type Product struct {
	Id    int     `json:"id" db:"product_id" `
	Name  string  `json:"name" db:"product_name"`
	Price float64 `json:"price" db:"product_price"`
}
