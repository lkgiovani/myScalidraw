package fileHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.fileUseCase.DeleteFile(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error deleting file"})
	}

	return c.JSON(fiber.Map{"message": "file deleted successfully"})
}
