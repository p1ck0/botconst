package repository

import (
	"context"

	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type Users interface {
	Create(ctx context.Context, user models.User) error
	GetByCredentials(ctx context.Context, email, password string) (models.User, error)
}

type Bots interface {
	Create(ctx context.Context, bot models.Bot) error
	GetByID(ctx context.Context, id string) (models.Bot, error)
	GetAll(ctx context.Context, userID string) ([]models.Bot, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, bot models.Bot) error
}

type Scenarios interface {
	Create(ctx context.Context, scenario models.Scenario) error
	GetByID(ctx context.Context, id string) (models.Scenario, error)
	GetAll(ctx context.Context, botID string) ([]models.Scenario, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, scenario models.Scenario) error
}

type Repositories struct {
	Users     Users
	Bots      Bots
	Scenarios Scenarios
}

func NewRepositories(db *mongo.Database) *Repositories {
	return &Repositories{
		Users:     mongodb.NewUsersRepo(db),
		Bots:      mongodb.NewBotsRepo(db),
		Scenarios: mongodb.NewScenariosRepo(db),
	}
}
