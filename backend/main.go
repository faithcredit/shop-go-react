package main

import (
	"context"
	"fmt"
	"gmrg/config"
	"gmrg/controllers"
	"gmrg/routes"
	"gmrg/services"
	"log"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

var (
	server      *gin.Engine
	ctx         context.Context
	mongoclient *mongo.Client
	redisclient *redis.Client
	ginCtx      *gin.Context

	userService         services.UserService
	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	authCollection      *mongo.Collection
	authService         services.AuthService
	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	// 👇 Add the Post Service, Controllers and Routes
	postService         services.PostService
	PostController      controllers.PostController
	postCollection      *mongo.Collection
	PostRouteController routes.PostRouteController

	// 👇 Add the Product Service, Controllers and Routes
	productCollection      *mongo.Collection
	productService         services.ProductService
	ProductController      controllers.ProductController
	ProductRouteController routes.ProductRouteController

	// 👇 Add the Review Service, Controllers and Routes
	reviewCollection      *mongo.Collection
	reviewService         services.ReviewService
	ReviewController      controllers.ReviewController
	ReviewRouteController routes.ReviewRouteController
)

func init() {
	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	ctx = context.TODO()

	// Connect to MongoDB
	mongoconn := options.Client().ApplyURI(config.DBUri)
	mongoclient, err := mongo.Connect(ctx, mongoconn)

	if err != nil {
		panic(err)
	}

	if err := mongoclient.Ping(ctx, readpref.Primary()); err != nil {
		panic(err)
	}

	fmt.Println("MongoDB successfully connected...")

	// Connect to Redis
	redisclient = redis.NewClient(&redis.Options{
		Addr: config.RedisUri,
	})

	if _, err := redisclient.Ping(ctx).Result(); err != nil {
		panic(err)
	}

	err = redisclient.Set(ctx, "test", "Welcome to Golang with Redis and MongoDB", 0).Err()
	if err != nil {
		panic(err)
	}

	fmt.Println("Redis client connected successfully...")

	// Collections
	authCollection = mongoclient.Database("golang_mongodb").Collection("users")
	userService = services.NewUserServiceImpl(authCollection, ctx)
	authService = services.NewAuthService(authCollection, ctx)
	AuthController = controllers.NewAuthController(authService, userService)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	UserController = controllers.NewUserController(userService)
	UserRouteController = routes.NewRouteUserController(UserController)
	// 👇 Add the Post Service, Controllers and Routes
	postCollection = mongoclient.Database("golang_mongodb").Collection("posts")
	postService = services.NewPostService(postCollection, ctx)
	PostController = controllers.NewPostController(postService)
	PostRouteController = routes.NewPostControllerRoute(PostController)
	//👇 Add the Product Service, Controllers and Routes
	productCollection = mongoclient.Database("golang_mongodb").Collection("products")
	productService = services.NewProductService(productCollection, ctx, ginCtx)
	ProductController = controllers.NewProductController(productService)
	ProductRouteController = routes.NewProductControllerRoute(ProductController)

	// 👇 Add the Review Service, Controllers and Routes
	reviewCollection = mongoclient.Database("golang_mongodb").Collection("reviews")
	reviewService = services.NewReviewService(reviewCollection, productCollection, ctx)
	ReviewController = controllers.NewReviewController(reviewService)
	ReviewRouteController = routes.NewReviewControllerRoute(ReviewController)

	server = gin.Default()
}

func main() {
	config, err := config.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load config", err)
	}

	defer mongoclient.Disconnect(ctx)

	value, err := redisclient.Get(ctx, "test").Result()

	if err == redis.Nil {
		fmt.Println("key: test does not exist")
	} else if err != nil {
		panic(err)
	}

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"http://localhost:5000", "http://localhost:3000"}
	corsConfig.AllowCredentials = true

	server.Use(cors.New(corsConfig))

	router := server.Group("/api")
	router.GET("/healthchecker", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"title": "success", "body": value})
	})

	AuthRouteController.AuthRoute(router, userService)
	UserRouteController.UserRoute(router, userService)
	// 👇 Evoke the PostRoute
	PostRouteController.PostRoute(router, userService)
	ReviewRouteController.ReviewRoute(router, userService)
	ProductRouteController.ProductRoute(router, userService)
	log.Fatal(server.Run(":" + config.Port))
}
