package handler

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/maxoov1/faq-api/pkg/service"
)

func (h *Handler) initWebHookRoutes(app *fiber.App) {
	webhook := app.Group("/webhook")

	webhook.Get("/facebook/:bot_id", h.facebook)
	webhook.Get("/testing/:hub.challenge.:species", h.testHub)
	webhook.Post("/:platform/:bot_id", h.handlerMessage)
}

func (h *Handler) facebook(ctx *fiber.Ctx) error {
	hubChallenge := strings.Split(ctx.String(), "&")[1]
	code := strings.Replace(hubChallenge, "hub.challenge=", "", 1)
	return ctx.SendString(code)
}

func (h *Handler) testHub(ctx *fiber.Ctx) error {
	fmt.Fprintf(ctx, "%s.%s\n", ctx.Params("genus"), ctx.Params("hub.verify_token"))
	return nil // prunus.persica
}

func (h *Handler) handlerMessage(ctx *fiber.Ctx) error {
	fmt.Println("GET!")
	botId := ctx.Params("bot_id")
	platform := ctx.Params("platform")
	switch platform {
	case "facebook":
		var message service.FacebookMessage
		err := ctx.BodyParser(&message)
		if err != nil {
			fmt.Println(err)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		for _, entry := range message.Entry {
			for _, msg := range entry.Messaging {
				if err != nil {
					fmt.Println(err)
					return ctx.SendStatus(fiber.StatusBadRequest)
				}
				err := h.service.WebHooks.HandlerMessage(ctx.Context(), service.Message{
					BotId:    botId,
					ChatId:   msg.Sender.ID,
					Platform: platform,
					Text:     msg.Message.Text,
				})
				if err != nil {
					return ctx.SendStatus(fiber.StatusBadRequest)
				}
			}
		}

	case "whatsapp":
		var response map[string][]map[string]interface{}
		err := json.Unmarshal(ctx.Body(), &response)
		if err != nil {
			fmt.Println(err)
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		for _, v := range response["messages"] {
			err := h.service.WebHooks.HandlerMessage(ctx.Context(), service.Message{
				BotId:    botId,
				ChatId:   v["author"].(string),
				Platform: platform,
				Text:     v["body"].(string),
			})
			if err != nil {
				return ctx.SendStatus(fiber.StatusBadRequest)
			}

		}
	case "telegram":
		webhook := service.TelegramMessage{}
		err := ctx.BodyParser(&webhook)
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}
		err = h.service.WebHooks.HandlerMessage(ctx.Context(), service.Message{
			BotId:    botId,
			ChatId:   fmt.Sprintf("%v", webhook.Message.From.ID),
			Platform: platform,
			Text:     webhook.Message.Text,
		})
		if err != nil {
			return ctx.SendStatus(fiber.StatusBadRequest)
		}

	default:
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
