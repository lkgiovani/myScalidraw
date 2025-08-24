package fileHandlers

import (
	"encoding/json"
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
	api.Put("/files/:id/rename", h.RenameFile)
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

	response := map[string]interface{}{
		"id":           file.ID,
		"name":         file.Name,
		"isFolder":     file.IsFolder,
		"parentId":     file.ParentID,
		"lastModified": file.LastModified,
		"path":         file.Path,
	}

	if !file.IsFolder {
		// Buscar o conteúdo original diretamente do repositório
		content, err := h.fileUseCase.GetFileContent(id)
		if err == nil && content != "" {
			response["content"] = content
		} else if file.Data != nil {
			// Fallback: usar Data se o conteúdo original não estiver disponível
			if contentBytes, err := json.Marshal(file.Data); err == nil {
				response["content"] = string(contentBytes)
			}
		}
	}

	return c.JSON(response)
}

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

func (h *FileHandler) DeleteFile(c *fiber.Ctx) error {
	id := c.Params("id")

	err := h.fileUseCase.DeleteFile(id)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error deleting file"})
	}

	return c.JSON(fiber.Map{"message": "file deleted successfully"})
}
