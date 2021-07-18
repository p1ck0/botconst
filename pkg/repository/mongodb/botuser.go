package mongodb

import (
	"context"
	"errors"

	"github.com/maxoov1/faq-api/pkg/database/mongodb"
	"github.com/maxoov1/faq-api/pkg/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type BotUserRepo struct {
	db *mongo.Collection
}

func NewBotUserRepo(db *mongo.Database) *BotUserRepo {
	return &BotUserRepo{
		db: db.Collection(botUserCollection),
	}
}

func (r *BotUserRepo) Create(ctx context.Context, botUser models.BotUser) error {
	_, err := r.db.InsertOne(ctx, botUser)
	if mongodb.IsDuplicate(err) {
		return errors.New("botuser already exists")
	}

	return err
}

func (r *BotUserRepo) GetByID(ctx context.Context, iduser, idbot, platform string) (models.BotUser, error) {
	var bot models.BotUser
	if err := r.db.FindOne(ctx, bson.M{
		"bot_user_id": iduser,
		"bot_id":      idbot,
		"platform":    platform,
	}).Decode(&bot); err != nil {
		return models.BotUser{}, err
	}

	return bot, nil
}

func (r *BotUserRepo) Update(ctx context.Context, botUser models.BotUser) error {
	_, err := r.db.UpdateOne(ctx, bson.M{
		"bot_user_id": botUser.BotUserId,
		"bot_id":      botUser.BotId,
		"platform":    botUser.Platform,
	}, bson.M{
		"$set": bson.M{
			"state": botUser.State,
		},
	})

	return err
}
