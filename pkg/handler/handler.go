package handler

import (
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/maxoov1/faq-api/pkg/config"
	"github.com/maxoov1/faq-api/pkg/service"
	"net/http"
)

type Handler struct {
	service *service.Services
}

func NewHandler(services *service.Services) *Handler {
	return &Handler{
		service: services,
	}
}

func (h *Handler) Init(cfg *config.Config) http.HandlerFunc {
	app := fiber.New()

	app.Use(
		logger.New(),
	)

	h.initUserRoutes(app)

	app.Use(
		jwtware.New(jwtware.Config{
			SigningKey: []byte(cfg.TokenManager.SigningKey),
		}),
	)

	h.initBotsRoutes(app)
	h.initScenariosRoutes(app)

	return adaptor.FiberApp(app)
}
