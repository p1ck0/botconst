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

type WebHookService struct {
	repoBots      repository.Bots
	repoScenarios repository.Scenarios
}

func NewWebHookService(repoBots repository.Bots, repoScenarios repository.Scenarios) *WebHookService {
	return &WebHookService{
		repoBots:      repoBots,
		repoScenarios: repoScenarios,
	}
}

func (s *WebHookService) HandlerMessage(ctx context.Context, m Message) error {
	actions, err := s.getActions(ctx, m.BotId, m.Text)
	if err != nil {
		return err
	}
	for _, action := range actions {
		fmt.Println(action.Type)
		if action.Type == "text" {
			fmt.Println(action.Value)
			s.sendMessage(ctx, Message{
				BotId:    m.BotId,
				ChatId:   m.ChatId,
				Platform: m.Platform,
				Text:     action.Value,
			})
		}
	}
	return nil
}

func (s *WebHookService) sendMessage(ctx context.Context, m Message) bool {
	bot, err := s.repoBots.GetByID(ctx, m.BotId)
	if err != nil {
		fmt.Println(err)
		return false
	}
	var sent bool
	fmt.Println("d", m.Platform)
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

func (s *WebHookService) getActions(ctx context.Context, botId string, text string) ([]models.Action, error) {
	scenarios, err := s.repoScenarios.GetAll(ctx, botId)
	if err != nil {
		fmt.Println(err)
		return []models.Action{}, err
	}
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

	return actions, nil
}
