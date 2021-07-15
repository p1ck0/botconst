package handler

import (
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
