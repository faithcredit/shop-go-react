package services

import "gmrg/models"

type BrandService interface {
	CreateBrand(*models.CreateBrandRequest) (*models.DBBrand, error)
	UpdateBrand(string, *models.UpdateBrand) (*models.DBBrand, error)
	FindBrandById(string) (*models.DBBrand, error)
	FindBrands(page int, limit int) ([]*models.DBBrand, error)
	DeleteBrand(string) error
}
