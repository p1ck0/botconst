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
	repoBotUser   repository.BotUser
}

func NewWebHookService(repoBots repository.Bots, repoScenarios repository.Scenarios, repoBotUser repository.BotUser) *WebHookService {
	return &WebHookService{
		repoBots:      repoBots,
		repoScenarios: repoScenarios,
		repoBotUser:   repoBotUser,
	}
}

func (s *WebHookService) HandlerMessage(ctx context.Context, m Message) error {
	botUser, err := s.repoBotUser.GetByID(ctx, m.ChatId, m.BotId, m.Platform)
	if err != nil {
		return err
	}
	actions, err := s.getActions(ctx, m, botUser)
	if err != nil {
		return err
	}
	return s.actionIter(ctx, actions, m, botUser)
}

func (s *WebHookService) actionIter(ctx context.Context, actions []models.Action, m Message, botUser models.BotUser) error {
	for _, action := range actions {
		switch action.Type {
		case "text":
			s.sendMessage(ctx, Message{
				BotId:    m.BotId,
				ChatId:   m.ChatId,
				Platform: m.Platform,
				Text:     action.Value,
			})
		case "redirect":
			botUser.State = ""
			if err := s.repoBotUser.Update(ctx, botUser); err != nil {
				return err
			}
			actions, err := s.getRedirectActions(ctx, m, botUser, action.Value)
			if err != nil {
				return err
			}
			s.actionIter(ctx, actions, m, botUser)
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

func (s *WebHookService) getActions(ctx context.Context, m Message, botUser models.BotUser) ([]models.Action, error) {
	var sent bool
	actions := []models.Action{}
	scenarios, err := s.repoScenarios.GetAll(ctx, m.BotId)
	if err != nil {
		fmt.Println(err)
		return []models.Action{}, err
	}
	for _, scenario := range scenarios {
		for _, trigger := range scenario.Triggers {
			if trigger == m.Text && botUser.State == "" {
				if err := s.repoBotUser.Create(ctx, models.BotUser{
					BotId:     m.BotId,
					BotUserId: m.ChatId,
					Platform:  m.Platform,
					State:     scenario.ID.Hex(),
				}); err != nil {
					err := s.repoBotUser.Update(ctx, models.BotUser{
						BotId:     m.BotId,
						BotUserId: m.ChatId,
						Platform:  m.Platform,
						State:     scenario.ID.Hex(),
					})
					return []models.Action{}, err
				}
				for _, dialog := range scenario.Dialogs {
					if dialog.IsMain {
						for _, action := range dialog.Actions {
							actions = append(actions, action)
							sent = true
						}
					}
				}
			} else if scenario.ID.Hex() == botUser.State {
				for _, dialog := range scenario.Dialogs {
					if dialog.Trigger == m.Text && !dialog.IsMain {
						for _, action := range dialog.Actions {
							actions = append(actions, action)
							sent = true
						}
					}
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

func (s *WebHookService) getRedirectActions(ctx context.Context, m Message, botUser models.BotUser, id string) ([]models.Action, error) {
	actions := []models.Action{}
	scenario, err := s.repoScenarios.GetByID(ctx, id)
	if err != nil {
		fmt.Println(err)
		return []models.Action{}, err
	}
	for _, trigger := range scenario.Triggers {
		if trigger == m.Text && botUser.State == "" {
			if err := s.repoBotUser.Create(ctx, models.BotUser{
				BotId:     m.BotId,
				BotUserId: m.ChatId,
				Platform:  m.Platform,
				State:     scenario.ID.Hex(),
			}); err != nil {
				err := s.repoBotUser.Update(ctx, models.BotUser{
					BotId:     m.BotId,
					BotUserId: m.ChatId,
					Platform:  m.Platform,
					State:     scenario.ID.Hex(),
				})
				return []models.Action{}, err
			}
			for _, dialog := range scenario.Dialogs {
				if dialog.IsMain {
					for _, action := range dialog.Actions {
						actions = append(actions, action)
					}
				}
			}
		} else if scenario.ID.Hex() == botUser.State {
			for _, dialog := range scenario.Dialogs {
				if dialog.Trigger == m.Text && !dialog.IsMain {
					for _, action := range dialog.Actions {
						actions = append(actions, action)
					}
				}
			}
		}
	}
	return actions, nil
}
