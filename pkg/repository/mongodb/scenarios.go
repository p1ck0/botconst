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

type ScenariosRepo struct {
	db *mongo.Collection
}

func NewScenariosRepo(db *mongo.Database) *ScenariosRepo {
	return &ScenariosRepo{
		db: db.Collection(scenariosCollection),
	}
}

func (r *ScenariosRepo) Create(ctx context.Context, scenario models.Scenario) error {
	_, err := r.db.InsertOne(ctx, scenario)
	if mongodb.IsDuplicate(err) {
		return errors.New("scenario already exists")
	}

	return err
}

func (r *ScenariosRepo) GetByID(ctx context.Context, id string) (models.Scenario, error) {
	var scenario models.Scenario
	idObj, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return models.Scenario{}, err
	}
	if err := r.db.FindOne(ctx, bson.M{
		"_id": idObj,
	}).Decode(&scenario); err != nil {
		return models.Scenario{}, err
	}

	return scenario, nil
}

func (r *ScenariosRepo) GetAll(ctx context.Context, botID string) ([]models.Scenario, error) {
	cursor, err := r.db.Find(ctx, bson.M{
		"botId": botID,
	})
	if err != nil {
		return nil, err
	}

	var scenarios []models.Scenario
	if err := cursor.All(ctx, &scenarios); err != nil {
		return nil, err
	}

	return scenarios, nil
}

func (r *ScenariosRepo) Delete(ctx context.Context, id string) error {
	_, err := r.db.DeleteOne(ctx, bson.M{
		"_id": id,
	})

	return err
}

func (r *ScenariosRepo) Update(ctx context.Context, scenario models.Scenario) error {
	query := bson.M{}

	if scenario.Name != "" {
		query["name"] = scenario.Name
	}

	if len(scenario.Triggers) != 0 {
		query["triggers"] = scenario.Triggers
	}

	if len(scenario.Dialogs) != 0 {
		query["actions"] = scenario.Dialogs
	}

	_, err := r.db.UpdateOne(ctx, bson.M{
		"_id": scenario.ID,
	}, bson.M{
		"$set": query,
	})

	return err
}
