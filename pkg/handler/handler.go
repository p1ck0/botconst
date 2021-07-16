package handler

import (
	"net/http"

	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/adaptor/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	jwtware "github.com/gofiber/jwt/v2"
	"github.com/maxoov1/faq-api/pkg/config"
	"github.com/maxoov1/faq-api/pkg/service"
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

	app.Get("/swagger/*", swagger.Handler) // default

	app.Get("/swagger/*", swagger.New(swagger.Config{ // custom
		URL:         "http://example.com/doc.json",
		DeepLinking: false,
	}))

	h.initWebHookRoutes(app)
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
