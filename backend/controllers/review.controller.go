package controllers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"gmrg/models"
	"gmrg/services"

	"github.com/gin-gonic/gin"
)

type ReviewController struct {
	reviewService services.ReviewService
}

func NewReviewController(reviewService services.ReviewService) ReviewController {
	return ReviewController{reviewService}
}

func (rc *ReviewController) CreateReview(ctx *gin.Context) {

	productId := ctx.Param("productId")
	fmt.Println("########CreateReview Controller############", productId)
	var review *models.CreateReviewRequest

	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	currentUser := ctx.MustGet("currentUser").(*models.DBResponse)
	review.UserID = currentUser.ID
	review.ProductID = productId
	newReview, err := rc.reviewService.CreateReview(review)

	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newReview})
}

func (rc *ReviewController) UpdateReview(ctx *gin.Context) {
	reviewId := ctx.Param("reviewId")

	var review *models.UpdateReview
	if err := ctx.ShouldBindJSON(&review); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updateReview, err := rc.reviewService.UpdateReview(reviewId, review)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updateReview})
}

func (rc *ReviewController) FindReviewById(ctx *gin.Context) {
	reviewId := ctx.Param("reviewId")

	review, err := rc.reviewService.FindReviewById(reviewId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": review})
}

func (rc *ReviewController) FindReviews(ctx *gin.Context) {
	var page = ctx.DefaultQuery("page", "1")
	var limit = ctx.DefaultQuery("limit", "10")

	intPage, err := strconv.Atoi(page)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	intLimit, err := strconv.Atoi(limit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	reviews, err := rc.reviewService.FindReviews(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(reviews), "data": reviews})
}

func (rc *ReviewController) DeleteReview(ctx *gin.Context) {
	reviewId := ctx.Param("reviewId")

	err := rc.reviewService.DeleteReview(reviewId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusNoContent, nil)
}
