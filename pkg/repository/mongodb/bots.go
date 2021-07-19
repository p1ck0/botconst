package mongodb

import (
	"context"
	"errors"

	"github.com/maxoov1/faq-api/pkg/database/mongodb"
	"github.com/maxoov1/faq-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type BotsRepo struct {
	db *mongo.Collection
}

func NewBotsRepo(db *mongo.Database) *BotsRepo {
	return &BotsRepo{
		db: db.Collection(botsCollection),
	}
}

func (r *BotsRepo) Create(ctx context.Context, bot models.Bot) error {
	_, err := r.db.InsertOne(ctx, bot)
	if mongodb.IsDuplicate(err) {
		return errors.New("bot already exists")
	}

	return err
}

func (r *BotsRepo) GetByID(ctx context.Context, id string) (models.Bot, error) {
	var bot models.Bot
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Bot{}, err
	}
	if err := r.db.FindOne(ctx, bson.M{
		"_id": objID,
	}).Decode(&bot); err != nil {
		return models.Bot{}, err
	}

	return bot, nil
}

func (r *BotsRepo) GetAll(ctx context.Context, userID string) ([]models.Bot, error) {
	cursor, err := r.db.Find(ctx, bson.M{
		"userId": userID,
	})
	if err != nil {
		return nil, err
	}

	var bots []models.Bot
	if err := cursor.All(ctx, &bots); err != nil {
		return nil, err
	}

	return bots, err
}

func (r *BotsRepo) Delete(ctx context.Context, id string) error {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	_, err = r.db.DeleteOne(ctx, bson.M{
		"_id": objID,
	})

	return err
}

func (r *BotsRepo) Update(ctx context.Context, bot models.Bot) error {
	query := bson.M{}

	if bot.UserID != "" {
		query["userId"] = bot.UserID
	}

	if bot.Name != "" {
		query["name"] = bot.Name
	}

	_, err := r.db.UpdateOne(ctx, bson.M{
		"_id": bot.ID,
	}, bson.M{
		"$set": query,
	})

	return err
}
