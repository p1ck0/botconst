package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxoov1/faq-api/pkg/auth"
	"github.com/maxoov1/faq-api/pkg/service"
)

func (h *Handler) initBotsRoutes(app *fiber.App) {
	bot := app.Group("/bot")

	bot.Post("/", h.botCreate)
	bot.Get("/:id", h.botGetByID)
	bot.Get("/", h.botGetAll)
	bot.Delete("/:id", h.botDelete)
	bot.Put("/", h.botUpdate)
}

// Create godoc
// @Summary create bot
// @Description  new bot
// @ID bot
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param bot body service.BotInput true "Bot Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /bot [post]
func (h *Handler) botCreate(ctx *fiber.Ctx) error {
	var bot service.BotInput
	if err := ctx.BodyParser(&bot); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}
	bot.UserID = auth.ParseToken(ctx)

	if err := h.service.Bots.Create(ctx.Context(), bot); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Create godoc
// @Summary Get bot by id
// @Description  get bot
// @ID get_bot
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Bot
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /bot/{id} [get]
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

// Create godoc
// @Summary Get all bot by id
// @Description  get all bot
// @ID get_all_bot
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Bot
// @Security ApiKeyAuth
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /bot [get]
func (h *Handler) botGetAll(ctx *fiber.Ctx) error {
	id := auth.ParseToken(ctx)

	bots, err := h.service.Bots.GetAll(ctx.Context(), id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": bots,
	})
}

// Create godoc
// @Summary DElete bot by id
// @Description  delete bot
// @ID delete_bot
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param id path string true "Bot Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /bot/{id} [delete]
func (h *Handler) botDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Bots.Delete(ctx.Context(), id); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Create godoc
// @Summary update bot
// @Description  up bot
// @ID update_bot
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param bot body service.BotInput true "Bot Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /bot [put]
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
