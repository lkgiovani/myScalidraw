package fileHandlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/files", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.GetFiles)
	api.Get("/files/:id", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.GetFileByID)
	api.Post("/files", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.CreateFile)
	api.Post("/files/upload", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.UploadFile)
	api.Put("/files/:id", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.SaveFile)
	api.Put("/files/:id/rename", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.RenameFile)
	api.Delete("/files/:id", h.authMiddleware.RequireSystemSetup(), h.authMiddleware.RequireAuth(), h.DeleteFile)
}
