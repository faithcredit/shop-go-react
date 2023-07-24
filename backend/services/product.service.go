package services

import "gmrg/models"

type ProductService interface {
	CreateProduct(*models.CreateProductRequest) (*models.DBProduct, error)
	UpdateProduct(string, *models.UpdateProduct) (*models.DBProduct, error)
	FindProductById(string) (*models.DBProduct, error)
	FindProducts(
		page int64,
		limit int64,
		category string,
		brand string,
		searchQuery string) ([]*models.DBProduct, error)
	DeleteProduct(string) error
	FindCategories(string) (models.Category, error)
	FindBrands(string) (models.Brand, error)
}
