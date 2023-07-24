package services

import (
	"context"
	"errors"
	"fmt"
	"gmrg/models"
	"gmrg/utils"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ReviewServieImpl struct {
	reviewCollection  *mongo.Collection
	productCollection *mongo.Collection
	ctx               context.Context
}

func NewReviewService(
	reviewCollection *mongo.Collection,
	productCollection *mongo.Collection,
	ctx context.Context,
) ReviewService {
	return &ReviewServieImpl{reviewCollection, productCollection, ctx}
}

func (r *ReviewServieImpl) CreateReview(review *models.CreateReviewRequest) (*models.DBReview, error) {

	obId, _ := primitive.ObjectIDFromHex(review.ProductID)
	fmt.Println("########CreateReview Service############", obId)

	productQuery := bson.M{"_id": obId}
	var product *models.DBProduct
	if err := r.productCollection.FindOne(r.ctx, productQuery).Decode(&product); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}
		return nil, err
	}

	review.CreateAt = time.Now()
	review.UpdatedAt = review.CreateAt
	review.Timestamps = true
	res, err := r.reviewCollection.InsertOne(r.ctx, review)

	if err != nil {
		if er, ok := err.(mongo.WriteException); ok && er.WriteErrors[0].Code == 11000 {
			return nil, errors.New("review with that name already exists")
		}
		return nil, err
	}

	opt := options.Index()
	opt.SetUnique(true)

	index := mongo.IndexModel{Keys: bson.M{"name": 1}, Options: opt}

	if _, err := r.reviewCollection.Indexes().CreateOne(r.ctx, index); err != nil {
		return nil, errors.New("could not create index for name")
	}

	var newReview *models.DBReview
	query := bson.M{"_id": res.InsertedID}
	if err = r.reviewCollection.FindOne(r.ctx, query).Decode(&newReview); err != nil {
		return nil, err
	}

	return newReview, nil
}

func (r *ReviewServieImpl) UpdateReview(id string, data *models.UpdateReview) (*models.DBReview, error) {
	doc, err := utils.ToDoc(data)
	if err != nil {
		return nil, err
	}

	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.D{{Key: "_id", Value: obId}}
	update := bson.D{{Key: "$set", Value: doc}}
	res := r.reviewCollection.FindOneAndUpdate(r.ctx, query, update, options.FindOneAndUpdate().SetReturnDocument(1))

	var updateReview *models.DBReview

	if err := res.Decode(&updateReview); err != nil {
		return nil, errors.New("no review with that Id exists")
	}

	return updateReview, nil
}

func (r *ReviewServieImpl) FindReviewById(id string) (*models.DBReview, error) {
	obId, _ := primitive.ObjectIDFromHex(id)

	query := bson.M{"_id": obId}

	var review *models.DBReview

	if err := r.reviewCollection.FindOne(r.ctx, query).Decode(&review); err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("no document with that Id exists")
		}

		return nil, err
	}

	return review, nil
}

func (r *ReviewServieImpl) FindReviews(page int, limit int) ([]*models.DBReview, error) {
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

	cursor, err := r.reviewCollection.Find(r.ctx, query, &opt)
	if err != nil {
		return nil, err
	}

	defer cursor.Close(r.ctx)

	var reviews []*models.DBReview

	for cursor.Next(r.ctx) {
		review := &models.DBReview{}
		err := cursor.Decode(review)

		if err != nil {
			return nil, err
		}

		reviews = append(reviews, review)
	}

	if err := cursor.Err(); err != nil {
		return nil, err
	}

	if len(reviews) == 0 {
		return []*models.DBReview{}, nil
	}

	return reviews, nil
}

func (r *ReviewServieImpl) DeleteReview(id string) error {
	obId, _ := primitive.ObjectIDFromHex(id)
	query := bson.M{"_id": obId}

	res, err := r.reviewCollection.DeleteOne(r.ctx, query)
	if err != nil {
		return err
	}

	if res.DeletedCount == 0 {
		return errors.New("no document with that Id exists")
	}

	return nil
}
