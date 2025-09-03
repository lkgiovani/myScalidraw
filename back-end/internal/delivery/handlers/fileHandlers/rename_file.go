package fileHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) RenameFile(c *fiber.Ctx) error {
	id := c.Params("id")

	var request struct {
		Name string `json:"name"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "error parsing request body"})
	}

	if request.Name == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "name is required"})
	}

	err := h.fileUseCase.RenameFile(id, request.Name)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error renaming file"})
	}

	updatedFile, err := h.fileUseCase.GetFileByID(id)
	if err != nil {

		return c.JSON(fiber.Map{"message": "file renamed successfully"})
	}

	return c.JSON(updatedFile)
}
