package fileHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) SaveFile(c *fiber.Ctx) error {
	id := c.Params("id")

	fileContent := c.Body()
	if len(fileContent) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "empty file content"})
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "content must be valid JSON",
			"details": err.Error(),
		})
	}

	if jsonData["type"] == nil {
		jsonData["type"] = "excalidraw"
	}
	if jsonData["version"] == nil {
		jsonData["version"] = 2
	}
	if jsonData["source"] == nil {
		jsonData["source"] = "https://excalidraw.com"
	}

	validatedContent, err := json.Marshal(jsonData)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error processing JSON"})
	}

	err = h.fileUseCase.SaveFile(id, string(validatedContent))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error saving file"})
	}

	updatedFile, err := h.fileUseCase.GetFileByID(id)
	if err != nil {

		return c.JSON(fiber.Map{"message": "file saved successfully"})
	}

	return c.JSON(updatedFile)
}
