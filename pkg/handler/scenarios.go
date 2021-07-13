package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/maxoov1/faq-api/pkg/service"
)

func (h *Handler) initScenariosRoutes(app *fiber.App) {
	scenario := app.Group("/scenario")

	scenario.Post("/", h.scenarioCreate)
	scenario.Get("/:id", h.scenarioGetByID)
	scenario.Get("/bot/:id", h.scenarioGetAll)
	scenario.Delete("/:id", h.scenarioDelete)
	scenario.Put("/", h.scenarioUpdate)
}

func (h *Handler) scenarioCreate(ctx *fiber.Ctx) error {
	var scenario service.ScenariosInput
	if err := ctx.BodyParser(&scenario); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Scenarios.Create(ctx.Context(), scenario); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) scenarioGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	scenario, err := h.service.Scenarios.GetByID(ctx.Context(), id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": scenario,
	})
}

func (h *Handler) scenarioGetAll(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	scenarios, err := h.service.Scenarios.GetAll(ctx.Context(), id)
	if err != nil {
		return ctx.SendStatus(fiber.StatusNotFound)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": scenarios,
	})
}

func (h *Handler) scenarioDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Scenarios.Delete(ctx.Context(), id); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) scenarioUpdate(ctx *fiber.Ctx) error {
	var scenario service.ScenariosInput
	if err := ctx.BodyParser(&scenario); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	if err := h.service.Scenarios.Update(ctx.Context(), scenario); err != nil {
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
