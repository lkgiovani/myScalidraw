package handlers

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

	app.Get("/api/excalidraw", h.GetExcalidraw)
	app.Get("/api/files", h.GetFiles)
	app.Get("/api/files/:id", h.GetFileByID)
}

func (h *FileHandler) GetExcalidraw(c *fiber.Ctx) error {
	data, exists := h.fileUseCase.GetSala("exemplo-salve")
	if !exists {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "arquivo não encontrado"})
	}

	c.Set("Content-Type", "application/json")
	return c.SendString(data)
}

func (h *FileHandler) GetFiles(c *fiber.Ctx) error {
	files := h.fileUseCase.GetFiles()
	return c.JSON(files)
}

func (h *FileHandler) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")
	file, err := h.fileUseCase.GetFileByID(id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao buscar arquivo"})
	}

	if file == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "arquivo não encontrado"})
	}

	return c.JSON(file)
}
