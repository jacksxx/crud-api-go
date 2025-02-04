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
