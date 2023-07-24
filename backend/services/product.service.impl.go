package services

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"gmrg/models"
	"gmrg/utils"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ProductServiceImpl struct {
	productCollection *mongo.Collection
	ctx               context.Context
	ginCtx            *gin.Context
}

func NewProductService(productCollection *mongo.Collection, ctx context.Context, ginCtx *gin.Context) ProductService {
	return &ProductServiceImpl{productCollection, ctx, ginCtx}
}

func (p *ProductServiceImpl) CreateProduct(product *models.CreateProductRequest) (*models.DBProduct, error) {
	product.CreateAt = time.Now()
	product.UpdatedAt = product.CreateAt
	res, err := p.productCollection.InsertOne(p.ctx, product)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("product with that name already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"name": 1}, Options: opt}

	if _, err := p.productCollection.Indexes().CreateOne(p.ctx, index); err != nil {
		return nil, errors.New("could not create index for name")
	}

	var newProduct *models.DBProduct
	query := bson.M{"_id": res.InsertedID}
	if err = p.productCollection.FindOne(p.ctx, query).Decode(&newProduct); err != nil {
		return nil, err
	}

	return newProduct, nil
}

func (p *ProductServiceImpl) UpdateProduct(id string, data *models.UpdateProduct) (*models.DBProduct, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := p.productCollection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedProduct *models.DBProduct

	if err := res.Decode(&updatedProduct); err != nil {
		return nil, errors.New("no product with that Id exists")
	}

	return updatedProduct, nil
}

func (p *ProductServiceImpl) FindProductById(id string) (*models.DBProduct, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var product *models.DBProduct

	if err := p.productCollection.FindOne(p.ctx, query).Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return product, nil
}
func (p *ProductServiceImpl) FindCategories(name string) (models.Category, error) {
	result, err := p.productCollection.Distinct(p.ctx, name, bson.M{})

	if err != nil {
		p.ginCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	distinctResult := models.Category{Categories: make([]string, len(result))}
	for i, v := range result {
		distinctResult.Categories[i] = v.(string)
	}
	return distinctResult, nil
}

func (p *ProductServiceImpl) FindBrands(name string) (models.Brand, error) {
	fmt.Println("###################################################", name)

	result, err := p.productCollection.Distinct(p.ctx, name, bson.M{})

	if err != nil {
		p.ginCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}

	distinctResult := models.Brand{Brands: make([]string, len(result))}
	for i, v := range result {

		distinctResult.Brands[i] = v.(string)
	}
	return distinctResult, nil
}

func (p *ProductServiceImpl) FindProducts(page int64,
	limit int64,
	category string,
	brand string,
	searchQuery string) ([]*models.DBProduct, error) {

	if page == 0 {
		page = 1
	}

	if limit == 0 {
		limit = 10
	}

	skip := (page - 1) * limit

	opt := options.FindOptions{}
	opt.SetLimit(int64(limit))
	opt.SetSkip(int64(skip))

	filter := bson.M{}

	if category != "" {
		filter["category"] = category
	}

	if brand != "" {
		filter["brand"] = brand
	}

	if searchQuery != "" {
		filter["name"] = bson.M{"$regex": searchQuery, "$options": "i"}
	}

	cursor, err := p.productCollection.Find(p.ctx, filter, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(p.ctx)

	var products []*models.DBProduct

	for cursor.Next(p.ctx) {
		product := &models.DBProduct{}
		err := cursor.Decode(product)

		if err != nil {
			return nil, err
		}

		products = append(products, product)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(products) == 0 {
		return []*models.DBProduct{}, nil
	}

	return products, nil
}

func (p *ProductServiceImpl) DeleteProduct(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := p.productCollection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
