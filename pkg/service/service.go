package service

import (
	"context"
	"time"

	"github.com/maxoov1/faq-api/pkg/auth"
	"github.com/maxoov1/faq-api/pkg/hash"
	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository"
)

type UserInputSignUp struct {
	Name     string
	Email    string
	Password string
}

type UserInputSignIn struct {
	Email    string
	Password string
}

type Users interface {
	SignUp(ctx context.Context, userInput UserInputSignUp) error
	SignIn(ctx context.Context, userInput UserInputSignIn) (string, error)
}

type BotInput struct {
	UserID              string
	Name                string
	Telegram            string
	WhatsAppID          string
	WhatsAppToken       string
	FacebookAppID       string
	FacebookAppSecret   string
	FacebookPageID      string
	FacebookAccessToken string
}

type Bots interface {
	Create(ctx context.Context, botInput BotInput) error
	GetByID(ctx context.Context, id string) (models.Bot, error)
	GetAll(ctx context.Context, userID string) ([]models.Bot, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, botInput BotInput) error
	HandlerMessage(ctx context.Context, m Message, scenarios []models.Scenario)
}

type ScenariosInput struct {
	BotID    string
	Name     string
	Triggers []string
	Actions  []models.Action
}

type Scenarios interface {
	Create(ctx context.Context, scenarioInput ScenariosInput) error
	GetByID(ctx context.Context, id string) (models.Scenario, error)
	GetAll(ctx context.Context, botID string) ([]models.Scenario, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, scenarioInput ScenariosInput) error
}

type Services struct {
	Users     Users
	Bots      Bots
	Scenarios Scenarios
}

func NewServices(r *repository.Repositories, hasher hash.Hasher, manager auth.TokenManager, tokenTTL time.Duration) *Services {
	return &Services{
		Users:     NewUsersService(r.Users, hasher, manager, tokenTTL),
		Bots:      NewBotsService(r.Bots),
		Scenarios: NewScenariosService(r.Scenarios),
	}
}

type FacebookMessage struct {
	Object string `json:"object"`
	Entry  []struct {
		ID        string `json:"id"`
		Time      int64  `json:"time"`
		Messaging []struct {
			Sender struct {
				ID string `json:"id"`
			} `json:"sender"`
			Recipient struct {
				ID string `json:"id"`
			} `json:"recipient"`
			Timestamp int64 `json:"timestamp"`
			Message   struct {
				Mid  string `json:"mid"`
				Text string `json:"text"`
			} `json:"message"`
		} `json:"messaging"`
	} `json:"entry"`
}

type TelegramMessage struct {
	UpdateID int `json:"update_id"`
	Message  struct {
		MessageID int `json:"message_id"`
		From      struct {
			ID           int    `json:"id"`
			IsBot        bool   `json:"is_bot"`
			FirstName    string `json:"first_name"`
			Username     string `json:"username"`
			LanguageCode string `json:"language_code"`
		} `json:"from"`
		Chat struct {
			ID        int    `json:"id"`
			FirstName string `json:"first_name"`
			Username  string `json:"username"`
			Type      string `json:"type"`
		} `json:"chat"`
		Date int    `json:"date"`
		Text string `json:"text"`
	} `json:"message"`
}

type Message struct {
	BotId    string
	ChatId   string
	Platform string
	Text     string
}
