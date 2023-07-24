package services

import (
	"context"
	"errors"
	"time"

	"gmrg/models"
	"gmrg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type BrandServiceImpl struct {
	brandCollection *mongo.Collection
	ctx             context.Context
}

func NewBrandService(brandCollection *mongo.Collection, ctx context.Context) BrandService {
	return &BrandServiceImpl{brandCollection, ctx}
}

func (p *BrandServiceImpl) CreateBrand(brand *models.CreateBrandRequest) (*models.DBBrand, error) {
	brand.CreateAt = time.Now()
	brand.UpdatedAt = brand.CreateAt
	res, err := p.brandCollection.InsertOne(p.ctx, brand)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("brand with that title already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: opt}

	if _, err := p.brandCollection.Indexes().CreateOne(p.ctx, index); err != nil {
		return nil, errors.New("could not create index for title")
	}

	var newBrand *models.DBBrand
	query := bson.M{"_id": res.InsertedID}
	if err = p.brandCollection.FindOne(p.ctx, query).Decode(&newBrand); err != nil {
		return nil, err
	}

	return newBrand, nil
}

func (p *BrandServiceImpl) UpdateBrand(id string, data *models.UpdateBrand) (*models.DBBrand, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := p.brandCollection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedBrand *models.DBBrand

	if err := res.Decode(&updatedBrand); err != nil {
		return nil, errors.New("no brand with that Id exists")
	}

	return updatedBrand, nil
}

func (p *BrandServiceImpl) FindBrandById(id string) (*models.DBBrand, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var brand *models.DBBrand

	if err := p.brandCollection.FindOne(p.ctx, query).Decode(&brand); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return brand, nil
}

func (p *BrandServiceImpl) FindBrands(page int, limit int) ([]*models.DBBrand, error) {
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

	query := bson.M{}

	cursor, err := p.brandCollection.Find(p.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(p.ctx)

	var brands []*models.DBBrand

	for cursor.Next(p.ctx) {
		brand := &models.DBBrand{}
		err := cursor.Decode(brand)

		if err != nil {
			return nil, err
		}

		brands = append(brands, brand)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(brands) == 0 {
		return []*models.DBBrand{}, nil
	}

	return brands, nil
}

func (p *BrandServiceImpl) DeleteBrand(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := p.brandCollection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
