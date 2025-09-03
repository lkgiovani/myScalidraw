package fileHandlers

import (
	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) GetFiles(c *fiber.Ctx) error {
	files := h.fileUseCase.GetFiles()
	return c.JSON(files)
}
