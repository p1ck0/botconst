package service

import (
	"context"

	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type BotsService struct {
	repo repository.Bots
}

func NewBotsService(repo repository.Bots) *BotsService {
	return &BotsService{repo: repo}
}

func (s *BotsService) Create(ctx context.Context, botInput BotInput) error {
	bot := models.Bot{
		ID:                  primitive.NewObjectID(),
		UserID:              botInput.UserID,
		Name:                botInput.Name,
		Telegram:            botInput.Telegram,
		WhatsAppID:          botInput.WhatsAppID,
		WhatsAppToken:       botInput.WhatsAppToken,
		FacebookAppID:       botInput.FacebookAppID,
		FacebookAppSecret:   botInput.FacebookAppSecret,
		FacebookPageID:      botInput.FacebookPageID,
		FacebookAccessToken: botInput.FacebookAccessToken,
	}

	if err := s.repo.Create(ctx, bot); err != nil {
		return err
	}
	setWebhooks(bot)
	return nil
}

func (s *BotsService) GetByID(ctx context.Context, id string) (models.Bot, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *BotsService) GetAll(ctx context.Context, userID string) ([]models.Bot, error) {
	return s.repo.GetAll(ctx, userID)
}

func (s *BotsService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func (s *BotsService) Update(ctx context.Context, botInput BotInput) error {
	bot := models.Bot{
		UserID:              botInput.UserID,
		Name:                botInput.Name,
		Telegram:            botInput.Telegram,
		WhatsAppID:          botInput.WhatsAppID,
		WhatsAppToken:       botInput.WhatsAppToken,
		FacebookAppID:       botInput.FacebookAppID,
		FacebookAppSecret:   botInput.FacebookAppSecret,
		FacebookPageID:      botInput.FacebookPageID,
		FacebookAccessToken: botInput.FacebookAccessToken,
	}

	return s.repo.Update(ctx, bot)
}
