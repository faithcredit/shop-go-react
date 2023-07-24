package services

import "gmrg/models"

type ReviewService interface {
	CreateReview(*models.CreateReviewRequest) (*models.DBReview, error)
	UpdateReview(string, *models.UpdateReview) (*models.DBReview, error)
	FindReviewById(string) (*models.DBReview, error)
	FindReviews(page int, limit int) ([]*models.DBReview, error)
	DeleteReview(string) error
}
