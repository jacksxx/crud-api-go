package service

import (
	"crud-api-go/arch/model"
	"crud-api-go/arch/repository"
)

type ProductService struct {
	//reposiroty
	ProductRepository repository.ProductRepository
}

func NewProductService(repository repository.ProductRepository) ProductService {
	return ProductService{
		ProductRepository: repository,
	}
}

func (p *ProductService) GetProducts() ([]model.Product, error) {
	return p.ProductRepository.GetProducts()
}

func (ps *ProductService) GetProductByID(id int) (*model.Product, error) {
	product, err := ps.ProductRepository.GetProductByID(id)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (ps *ProductService) CreateProducts(product model.Product) (model.Product, error) {

	productId, err := ps.ProductRepository.CreateProducts(product)
	if err != nil {
		return model.Product{}, err
	}
	product.Id = productId

	return product, nil
}
