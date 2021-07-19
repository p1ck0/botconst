package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/imroc/req"
	"github.com/maxoov1/faq-api/pkg/auth"
	"github.com/maxoov1/faq-api/pkg/config"
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
}

type ScenariosInput struct {
	ID       string
	BotID    string
	Name     string
	Triggers []string
	Dialogs  []models.Dialog
}

type Scenarios interface {
	Create(ctx context.Context, scenarioInput ScenariosInput) error
	GetByID(ctx context.Context, id string) (models.Scenario, error)
	GetAll(ctx context.Context, botID string) ([]models.Scenario, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, scenarioInput ScenariosInput) error
}

type WebHooks interface {
	HandlerMessage(ctx context.Context, m Message) error
	sendMessage(ctx context.Context, m Message) bool
	getActions(ctx context.Context, m Message, botUser models.BotUser) ([]models.Action, error)
	getRedirectActions(ctx context.Context, m Message, botUser models.BotUser, id string) ([]models.Action, error)
	actionIter(ctx context.Context, actions []models.Action, m Message, botUser models.BotUser) error
}

type Services struct {
	Users     Users
	Bots      Bots
	Scenarios Scenarios
	WebHooks  WebHooks
}

func NewServices(r *repository.Repositories, hasher hash.Hasher, manager auth.TokenManager, tokenTTL time.Duration) *Services {
	return &Services{
		Users:     NewUsersService(r.Users, hasher, manager, tokenTTL),
		Bots:      NewBotsService(r.Bots),
		Scenarios: NewScenariosService(r.Scenarios),
		WebHooks:  NewWebHookService(r.Bots, r.Scenarios, r.BotUser),
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

func setWebhooks(b models.Bot) {
	if b.Telegram != "" {
		WebhookURL := fmt.Sprintf("%s/%s/%s", config.Host, "telegram", b.ID.Hex())
		fmt.Println(WebhookURL)
		url := fmt.Sprintf("https://api.telegram.org/bot%s/setWebhook?url=%s", b.Telegram, WebhookURL)
		r, err := req.New().Get(url)
		var response map[string]interface{}
		err = r.ToJSON(&response)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(&response)
	}
	if b.WhatsAppToken != "" || b.WhatsAppID != "" {
		WebhookURL := fmt.Sprintf("%s/%s/%s", config.Host, "whatsapp", b.ID.Hex())
		url := fmt.Sprintf("https://api.chat-api.com/instance%s/webhook", b.WhatsAppID)
		r, err := req.Post(url, req.Param{"webhookUrl": WebhookURL}, req.QueryParam{
			"token": b.WhatsAppToken,
		})
		if err != nil {
			fmt.Println(err)
		}
		var response map[string]interface{}
		err = r.ToJSON(&response)
		if err != nil {
			fmt.Println(err)
		}
		if response["set"] == true {
			fmt.Println("Webhook set", response)
		} else {
			fmt.Println("Webhook not set")
			fmt.Println(response)
		}
	}
}
