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

type ProductController struct {
	productService services.ProductService
}

func NewProductController(productService services.ProductService) ProductController {
	return ProductController{productService}
}

func (pc *ProductController) CreateProduct(ctx *gin.Context) {
	var product *models.CreateProductRequest

	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}

	newProduct, err := pc.productService.CreateProduct(product)

	if err != nil {
		if strings.Contains(err.Error(), "name already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newProduct})
}

func (pc *ProductController) UpdateProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	var product *models.UpdateProduct
	if err := ctx.ShouldBindJSON(&product); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedProduct, err := pc.productService.UpdateProduct(productId, product)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedProduct})
}

func (pc *ProductController) FindProductById(ctx *gin.Context) {
	productId := ctx.Param("productId")

	product, err := pc.productService.FindProductById(productId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": product})
}

func (pc *ProductController) FindProducts(ctx *gin.Context) {

	fmt.Println("##########################react request################")
	page := ctx.DefaultQuery("page", "1")
	pageSize := ctx.DefaultQuery("limit", "9")
	category := ctx.DefaultQuery("category", "")
	brand := ctx.DefaultQuery("brand", "")
	query := ctx.DefaultQuery("query", "")

	initPage, err := strconv.ParseInt(page, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	// initPageSize, err := strconv.Atoi(pageSize)
	initPageSize, err := strconv.ParseInt(pageSize, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	products, err := pc.productService.FindProducts(
		initPage,
		initPageSize,
		category,
		brand,
		query,
	)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	brands, err := pc.productService.FindBrands("brand")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	categories, err := pc.productService.FindCategories("category")
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	stringLen := strconv.Itoa(len(products))
	f, err := strconv.ParseInt(stringLen, 10, 64)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid number"})
		return
	}
	result := (f / initPageSize)

	ctx.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"countProducts": len(products),
		"productDocs":   products,
		"categories":    categories,
		"brands":        brands,
		"page":          initPage,
		"pages":         result,
	})
}

func (pc *ProductController) DeleteProduct(ctx *gin.Context) {
	productId := ctx.Param("productId")

	err := pc.productService.DeleteProduct(productId)

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
