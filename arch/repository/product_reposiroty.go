package repository

import (
	"crud-api-go/arch/model"
	"database/sql"
	"fmt"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := `SELECT products_id, products_name, products_price FROM prod.products`
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}
	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(&productObj.Id, &productObj.Name, &productObj.Price)
		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}
		productList = append(productList, productObj)
	}
	rows.Close()
	return productList, nil
}

func (pr *ProductRepository) GetProductByID(id int) (*model.Product, error) {
	query, err := pr.connection.Prepare("SELECT * FROM prod.products WHERE products_id = $1")
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	var product model.Product
	err = query.QueryRow(id).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	query.Close()
	return &product, nil
}

func (pr *ProductRepository) CreateProducts(product model.Product) (int, error) {
	var Id int
	query, err := pr.connection.Prepare("INSERT INTO prod.products" + "(products_name, products_price)" + "VALUES ($1,$2) RETURNING products_id")

	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	defer query.Close()

	err = query.QueryRow(product.Name, product.Price).Scan(&Id)
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return Id, nil
}
