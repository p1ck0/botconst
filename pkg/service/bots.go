package service

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"github.com/imroc/req"
	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository"
	messenger "github.com/mileusna/facebook-messenger"
)

type BotsService struct {
	repo repository.Bots
}

func NewBotsService(repo repository.Bots) *BotsService {
	return &BotsService{repo: repo}
}

func (s *BotsService) Create(ctx context.Context, botInput BotInput) error {
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

	return s.repo.Create(ctx, bot)
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

func (s *BotsService) HandlerMessage(ctx context.Context, m Message, scenarios []models.Scenario) {
	actions := GetActions(m.BotId, m.Text, scenarios)
	for _, action := range actions {
		if action.Type == "text" {
			s.sendMessage(ctx, Message{
				BotId:    m.BotId,
				ChatId:   m.ChatId,
				Platform: m.Platform,
				Text:     action.Value,
			})
		}
	}
}

func (s *BotsService) sendMessage(ctx context.Context, m Message) bool {
	bot, err := s.repo.GetByID(ctx, m.BotId)
	if err != nil {
		return false
	}
	var sent bool
	if m.Platform == "facebook" {
		messenger := &messenger.Messenger{
			AccessToken: bot.FacebookAccessToken,
			VerifyToken: bot.FacebookAppSecret,
			PageID:      bot.FacebookPageID,
		}
		senderId, _ := strconv.Atoi(m.ChatId)
		messenger.SendTextMessage(int64(senderId), m.Text)
	}
	if m.Platform == "telegram" {
		url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage?chat_id=%s&text=%s", bot.Telegram, m.ChatId, m.Text)
		r, err := req.New().Get(url)
		if err != nil {
			log.Fatal(err)
		}
		var body map[string]interface{}
		r.ToJSON(&body)
		if body["ok"] == true {
			sent = true
		}
	}

	if m.Platform == "whatsapp" {
		url := fmt.Sprintf("https://api.chat-api.com/instance%s/sendMessage", bot.WhatsAppID)
		r, err := req.Post(url, req.Param{"chatId": m.ChatId, "body": m.Text}, req.QueryParam{"token": bot.WhatsAppToken})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(r.Response().StatusCode)
		var response map[string]interface{}
		r.ToJSON(response)
		fmt.Println(response)
	}

	return sent
}

func GetActions(botId string, text string, scenarios []models.Scenario) []models.Action {
	var sent bool
	actions := []models.Action{}
	for _, s := range scenarios {
		for _, trigger := range s.Triggers {
			if trigger == text {
				for _, action := range s.Actions {
					actions = append(actions, action)
					sent = true
				}
			}
		}
	}

	if sent == false {
		actions = append(actions, models.Action{
			Type:  "text",
			Value: "Я пока не обучен, чтобы отвечать на такие вопросы.",
		})
	}
	return actions
}
