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

// Create godoc
// @Summary create scenario
// @Description  new scenario
// @ID scenario
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param scenario  body service.ScenariosInput true "scenario Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /scenario [post]
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

// Create godoc
// @Summary Get scenario by id
// @Description  get scenario
// @ID get_scenario
// @Accept  json
// @Produce  json
// @Success 200 {object} models.Scenario
// @Security ApiKeyAuth
// @Param id path string true "scenario Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /scenario/{id} [get]
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

// Create godoc
// @Summary Get all scenarios by id
// @Description  get all scenarios
// @ID get_all_scenarios
// @Accept  json
// @Produce  json
// @Success 200 {object} []models.Scenario
// @Security ApiKeyAuth
// @Param bot_id path string true "scenario Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /scenario/bot/{id} [get]
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

// Create godoc
// @Summary DElete scenario by id
// @Description  delete scenario
// @ID delete_scenario
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param id path string true "scenario Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /scenario/{id} [delete]
func (h *Handler) scenarioDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")

	if err := h.service.Scenarios.Delete(ctx.Context(), id); err != nil {
		return ctx.SendStatus(fiber.StatusBadRequest)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Create godoc
// @Summary update scenario
// @Description  up scenario
// @ID update_scenario
// @Accept  json
// @Produce  json
// @Success 200
// @Security ApiKeyAuth
// @Param scenario body service.ScenariosInput true "scenario Data"
// @Param Authorization header string true "Insert your access token" default(Bearer <Add access token here>)
// @Router /scenario [put]
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
