package handler

import (
	"encoding/json"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maxoov1/faq-api/pkg/service"
)

func (h *Handler) initBotsRoutes(app *fiber.App) {
	bot := app.Group("/bot")

	bot.Post("/", h.botCreate)
	bot.Get("/:id", h.botGetByID)
	bot.Get("/user/:id", h.botGetAll)
	bot.Delete("/:id", h.botDelete)
	bot.Put("/", h.botUpdate)
	bot.Post("/:platform/:bot_id", h.botSendMessage)
}

func (h *Handler) botCreate(ctx *fiber.Ctx) error {
	var bot service.BotInput
	if err := ctx.BodyParser(&bot); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Bots.Create(ctx.Context(), bot); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) botGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	bot, err := h.service.Bots.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": bot,
	})
}

func (h *Handler) botGetAll(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	bots, err := h.service.Bots.GetAll(ctx.Context(), id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": bots,
	})
}

func (h *Handler) botDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Bots.Delete(ctx.Context(), id); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) botUpdate(ctx *fiber.Ctx) error {
	var bot service.BotInput
	if err := ctx.BodyParser(&bot); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Bots.Update(ctx.Context(), bot); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) botSendMessage(ctx *fiber.Ctx) error {
	botId := ctx.Params("bot_id")
	platform := ctx.Params("platform")
	scenarois, err := h.service.Scenarios.GetAll(ctx.Context(), botId)
	if err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	switch platform {
	case "facebook":
		var message service.FacebookMessage
		err := ctx.BodyParser(&message)
		if err != nil {
			fmt.Println(err)
		}
		for _, entry := range message.Entry {
			for _, msg := range entry.Messaging {
				if err != nil {
					fmt.Println(err)
				}
				h.service.Bots.HandlerMessage(ctx.Context(), service.Message{
					BotId:    botId,
					ChatId:   msg.Sender.ID,
					Platform: platform,
					Text:     msg.Message.Text,
				}, scenarois)
			}
		}

	case "whatsapp":
		var response map[string][]map[string]interface{}
		err := json.Unmarshal(ctx.Body(), &response)
		if err != nil {
			fmt.Println(err)
		}
		for _, v := range response["messages"] {
			h.service.Bots.HandlerMessage(ctx.Context(), service.Message{
				BotId:    botId,
				ChatId:   v["author"].(string),
				Platform: platform,
				Text:     v["body"].(string),
			}, scenarois)

		}
	case "telegram":
		webhook := service.TelegramMessage{}
		err := ctx.BodyParser(&webhook)
		if err != nil {
			fmt.Println(err)
		}
		h.service.Bots.HandlerMessage(ctx.Context(), service.Message{
			BotId:    botId,
			ChatId:   fmt.Sprintf("%v", webhook.Message.From.ID),
			Platform: platform,
			Text:     webhook.Message.Text,
		}, scenarois)

	default:
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
