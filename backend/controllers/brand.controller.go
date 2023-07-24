package controllers

import (
	"net/http"
	"strconv"
	"strings"

	"gmrg/models"
	"gmrg/services"

	"github.com/gin-gonic/gin"
)

type BrandController struct {
	brandService services.BrandService
}

func NewBrandController(brandService services.BrandService) BrandController {
	return BrandController{brandService}
}

func (pc *BrandController) CreateBrand(ctx *gin.Context) {
	var brand *models.CreateBrandRequest

	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadRequest, err.Error())
		return
	}
	newBrand, err := pc.brandService.CreateBrand(brand)

	if err != nil {
		if strings.Contains(err.Error(), "title already exists") {
			ctx.JSON(http.StatusConflict, gin.H{"status": "fail", "message": err.Error()})
			return
		}

		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "success", "data": newBrand})
}

func (pc *BrandController) UpdateBrand(ctx *gin.Context) {
	brandId := ctx.Param("brandId")

	var brand *models.UpdateBrand
	if err := ctx.ShouldBindJSON(&brand); err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	updatedBrand, err := pc.brandService.UpdateBrand(brandId, brand)
	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": updatedBrand})
}

func (pc *BrandController) FindBrandById(ctx *gin.Context) {
	brandId := ctx.Param("brandId")

	brand, err := pc.brandService.FindBrandById(brandId)

	if err != nil {
		if strings.Contains(err.Error(), "Id exists") {
			ctx.JSON(http.StatusNotFound, gin.H{"status": "fail", "message": err.Error()})
			return
		}
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "data": brand})
}

func (pc *BrandController) FindBrands(ctx *gin.Context) {
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

	brands, err := pc.brandService.FindBrands(intPage, intLimit)
	if err != nil {
		ctx.JSON(http.StatusBadGateway, gin.H{"status": "fail", "message": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"status": "success", "results": len(brands), "data": brands})
}

func (pc *BrandController) DeleteBrand(ctx *gin.Context) {
	brandId := ctx.Param("brandId")

	err := pc.brandService.DeleteBrand(brandId)

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
