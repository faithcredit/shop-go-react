package services

import "gmrg/models"

type CategoryService interface {
	CreateCategory(*models.CreateCategoryRequest) (*models.DBCategory, error)
	UpdateCategory(string, *models.UpdateCategory) (*models.DBCategory, error)
	FindCategoryById(string) (*models.DBCategory, error)
	FindCategorys(page int, limit int) ([]*models.DBCategory, error)
	DeleteCategory(string) error
}
