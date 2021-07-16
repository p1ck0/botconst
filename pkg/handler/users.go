package handler

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/maxoov1/faq-api/pkg/service"
)

type userInputSignUp struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userInputSignIn struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type Token struct {
	Token string `json:"token"`
}

func (h *Handler) initUserRoutes(app *fiber.App) {
	user := app.Group("/user")

	user.Post("/sign-up", h.userSignUp)
	user.Post("/sign-in", h.userSignIn)
}

// Create godoc
// @Summary Create a account
// @Description Create new account
// @ID create-user
// @Accept  json
// @Produce  json
// @Param user body userInputSignUp true "User Data"
// @Success 200
// @Router /user/sign-up [post]
func (h *Handler) userSignUp(ctx *fiber.Ctx) error {
	var input userInputSignUp
	if err := ctx.BodyParser(&input); err != nil {
		fmt.Println(err)
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	if err := h.service.Users.SignUp(ctx.Context(), service.UserInputSignUp{
		Name:     input.Name,
		Email:    input.Email,
		Password: input.Password,
	}); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.SendStatus(fiber.StatusOK)
}

// Create godoc
// @Summary sign-in a account
// @Description sign-in new account
// @ID sign-in
// @Accept  json
// @Produce  json
// @Param user body userInputSignIn true "User Data"
// @Success 200 {object} Token
// @Router /user/sign-in [post]
func (h *Handler) userSignIn(ctx *fiber.Ctx) error {
	var input userInputSignIn
	if err := ctx.BodyParser(&input); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	token, err := h.service.Users.SignIn(ctx.Context(), service.UserInputSignIn{
		Email:    input.Email,
		Password: input.Password,
	})
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"token": token,
	})
}
