package fileHandlers

import (
	"encoding/json"
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) GetFileByID(c *fiber.Ctx) error {
	id := c.Params("id")

	file, err := h.fileUseCase.GetFileByID(id)

	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error fetching file"})
	}

	if file == nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{"error": "file not found"})
	}

	response := map[string]interface{}{
		"id":           file.ID,
		"name":         file.Name,
		"isFolder":     file.IsFolder,
		"parentId":     file.ParentID,
		"lastModified": file.LastModified,
		"path":         file.Path,
	}

	if !file.IsFolder {
		content, err := h.fileUseCase.GetFileContent(id)
		if err == nil && content != "" {
			response["content"] = content
		} else if file.Data != nil {
			if contentBytes, err := json.Marshal(file.Data); err == nil {
				response["content"] = string(contentBytes)
			}
		}
	}

	return c.JSON(response)
}
