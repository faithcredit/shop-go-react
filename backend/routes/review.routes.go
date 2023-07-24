package routes

import (
	"gmrg/controllers"
	"gmrg/middleware"
	"gmrg/services"

	"github.com/gin-gonic/gin"
)

type ReviewRouteController struct {
	reviewController controllers.ReviewController
}

func NewReviewControllerRoute(reviewController controllers.ReviewController) ReviewRouteController {
	return ReviewRouteController{reviewController}
}

func (rrc *ReviewRouteController) ReviewRoute(rg *gin.RouterGroup, userService services.UserService) {
	protected := rg.Group("/:productId/auth/reviews")
	protected.Use(middleware.DeserializeUser(userService))
	router := rg.Group("/reviews")

	router.GET("/", rrc.reviewController.FindReviews)
	router.GET("/:reviewId", rrc.reviewController.FindReviewById)
	protected.POST("/", rrc.reviewController.CreateReview)
	router.PATCH("/:reviewId", rrc.reviewController.UpdateReview)
	router.DELETE("/:reviewId", rrc.reviewController.DeleteReview)
}
