package userHandlers

import (
	"github.com/gofiber/fiber/v2"

	"myScalidraw/internal/delivery/middleware"
)

func (h *UserHandler) RegisterRoutes(app *fiber.App, authMiddleware *middleware.AuthMiddleware) {
	api := app.Group("/api/auth")

	api.Post("/create-first-user", h.CreateFirstUser)
	api.Get("/has-users", h.HasUsers)

	api.Post("/login", authMiddleware.RequireSystemSetup(), h.Login)
	api.Post("/logout", authMiddleware.RequireSystemSetup(), h.Logout)
	api.Get("/me", authMiddleware.RequireSystemSetup(), h.GetCurrentUser)
	api.Get("/users", authMiddleware.RequireSystemSetup(), h.GetUsers)
}
