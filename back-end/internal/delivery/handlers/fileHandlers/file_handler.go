package fileHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"myScalidraw/internal/domain/useCase/file"
)

type FileHandler struct {
	fileUseCase *file.FileUseCase
}

func NewFileHandler(fileUseCase *file.FileUseCase) *FileHandler {
	return &FileHandler{
		fileUseCase: fileUseCase,
	}
}

func (h *FileHandler) RegisterRoutes(app *fiber.App) {
	api := app.Group("/api")

	api.Get("/files", h.GetFiles)
	api.Get("/files/:id", h.GetFileByID)
	api.Post("/files", h.CreateFile)
	api.Post("/files/upload", h.UploadFile)
	api.Put("/files/:id", h.SaveFile)
	api.Delete("/files/:id", h.DeleteFile)
}

func (h *FileHandler) GetFiles(c *fiber.Ctx) error {
	files := h.fileUseCase.GetFiles()
	return c.JSON(files)
}

func (h *FileHandler) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := h.fileUseCase.GetFileByID(id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching file"})
	}

	if file == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
	}

	return c.JSON(file)
}

func (h *FileHandler) SaveFile(c *fiber.Ctx) error {
	id := c.Params("id")

	var request struct {
		Content string `json:"content"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "error parsing request body"})
	}

	err := h.fileUseCase.SaveFile(id, request.Content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error saving file"})
	}

	return c.JSON(fiber.Map{"message": "file saved successfully"})
}

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.fileUseCase.DeleteFile(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error deleting file"})
	}

	return c.JSON(fiber.Map{"message": "file deleted successfully"})
}
