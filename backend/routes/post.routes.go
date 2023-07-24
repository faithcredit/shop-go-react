package routes

import (
	"gmrg/controllers"
	"gmrg/middleware"
	"gmrg/services"

	"github.com/gin-gonic/gin"
)

type PostRouteController struct {
	postController controllers.PostController
}

func NewPostControllerRoute(postController controllers.PostController) PostRouteController {
	return PostRouteController{postController}
}

func (r *PostRouteController) PostRoute(rg *gin.RouterGroup, userService services.UserService) {

	protected := rg.Group("/auth/posts")
	protected.Use(middleware.DeserializeUser(userService))
	router := rg.Group("/posts")

	router.GET("/", r.postController.FindPosts)
	router.GET("/:postId", r.postController.FindPostById)
	protected.POST("/", r.postController.CreatePost)
	protected.PATCH("/:postId", r.postController.UpdatePost)
	protected.DELETE("/:postId", r.postController.DeletePost)
}
