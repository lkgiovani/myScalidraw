package fileHandlers

import (
	"encoding/json"
	"myScalidraw/internal/domain/models/file_model"
	"myScalidraw/pkg/uuid"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	fileContent := c.Body()
	if len(fileContent) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "empty or not found file"})
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

	fileName := c.Get("X-File-Name")
	if fileName == "" {
		fileName = "Untitled-" + time.Now().Format("2006-01-02-1504") + ".excalidraw"
	} else if !strings.HasSuffix(fileName, ".excalidraw") {
		fileName = strings.TrimSuffix(fileName, ".json") + ".excalidraw"
	}

	parentID := c.Query("parentId")
	if parentID == "" {
		parentID = "drafts"
	}

	contentType := "application/vnd.excalidraw+json"

	fileID, err := uuid.GenerateUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error generating UUID"})
	}

	metadata := &file_model.FileMetadata{
		ID:           fileID,
		ParentID:     parentID,
		Name:         fileName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsFolder:     false,
		Size:         int64(len(validatedContent)),
		ContentType:  contentType,
		LastModified: time.Now(),
	}

	err = h.fileUseCase.CreateFile(metadata, validatedContent)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error uploading file"})
	}

	response := map[string]interface{}{
		"id":           metadata.ID,
		"name":         metadata.Name,
		"path":         "/" + metadata.Name,
		"size":         metadata.Size,
		"modified":     metadata.LastModified.Format(time.RFC3339),
		"lastModified": metadata.LastModified.Unix() * 1000,
		"isFolder":     metadata.IsFolder,
		"parentId":     metadata.ParentID,
	}

	return c.JSON(response)
}
