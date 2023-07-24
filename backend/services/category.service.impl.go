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

type CategoryServiceImpl struct {
	categoryCollection *mongo.Collection
	productCollection  *mongo.Collection
	ctx                context.Context
}

func NewCategoryService(categoryCollection *mongo.Collection, productCollection *mongo.Collection, ctx context.Context) CategoryService {
	return &CategoryServiceImpl{categoryCollection, productCollection, ctx}
}

func (p *CategoryServiceImpl) CreateCategory(category *models.CreateCategoryRequest) (*models.DBCategory, error) {
	category.CreateAt = time.Now()
	category.UpdatedAt = category.CreateAt
	res, err := p.categoryCollection.InsertOne(p.ctx, category)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("category with that title already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"title": 1}, Options: opt}

	if _, err := p.categoryCollection.Indexes().CreateOne(p.ctx, index); err != nil {
		return nil, errors.New("could not create index for title")
	}

	var newCategory *models.DBCategory
	query := bson.M{"_id": res.InsertedID}
	if err = p.categoryCollection.FindOne(p.ctx, query).Decode(&newCategory); err != nil {
		return nil, err
	}

	return newCategory, nil
}

func (p *CategoryServiceImpl) UpdateCategory(id string, data *models.UpdateCategory) (*models.DBCategory, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := p.categoryCollection.FindOneAndUpdate(p.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updatedCategory *models.DBCategory

	if err := res.Decode(&updatedCategory); err != nil {
		return nil, errors.New("no category with that Id exists")
	}

	return updatedCategory, nil
}

func (p *CategoryServiceImpl) FindCategoryById(id string) (*models.DBCategory, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var category *models.DBCategory

	if err := p.categoryCollection.FindOne(p.ctx, query).Decode(&category); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return category, nil
}

func (p *CategoryServiceImpl) FindCategorys(page int, limit int) ([]*models.DBCategory, error) {
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

	cursor, err := p.categoryCollection.Find(p.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(p.ctx)

	var categorys []*models.DBCategory

	for cursor.Next(p.ctx) {
		category := &models.DBCategory{}
		err := cursor.Decode(category)

		if err != nil {
			return nil, err
		}

		categorys = append(categorys, category)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(categorys) == 0 {
		return []*models.DBCategory{}, nil
	}

	return categorys, nil
}

func (p *CategoryServiceImpl) DeleteCategory(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := p.categoryCollection.DeleteOne(p.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
