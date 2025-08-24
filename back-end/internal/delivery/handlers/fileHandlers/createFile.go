package fileHandlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"myScalidraw/internal/domain/models"
	"myScalidraw/pkg/uuid"

	"github.com/gofiber/fiber/v2"
)

type CreateFileRequest struct {
	Name        string `json:"name"`
	ParentID    string `json:"parentId"`
	ParentPath  string `json:"parentPath"`
	IsFolder    bool   `json:"isFolder"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

func (h *FileHandler) CreateFile(c *fiber.Ctx) error {
	// Log do body raw para debug
	rawBody := c.Body()
	fmt.Printf("CreateFile - Raw body: %s\n", string(rawBody))

	var params CreateFileRequest
	if err := c.BodyParser(&params); err != nil {
		fmt.Printf("CreateFile - Error parsing body: %v\n", err)
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "error parsing request body"})
	}

	fmt.Printf("CreateFile - Parsed params: %+v\n", params)

	fileID, err := uuid.GenerateUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error generating UUID"})
	}

	fileName := params.Name

	// Se não especificado, assume que é um arquivo (não pasta)
	isFolder := params.IsFolder
	if params.Content != "" {
		isFolder = false // Se tem conteúdo, definitivamente é um arquivo
	}

	if !isFolder && !strings.HasSuffix(fileName, ".excalidraw") {
		fileName = strings.TrimSuffix(fileName, ".json") + ".excalidraw"
	}

	metadata := &models.FileMetadata{
		ID:           fileID,
		ParentID:     params.ParentID,
		Name:         fileName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsFolder:     isFolder,
		Size:         0,
		ContentType:  "application/vnd.excalidraw+json",
		LastModified: time.Now(),
	}

	var content []byte
	if !isFolder {
		if params.Content != "" {
			content = []byte(params.Content)
		} else {
			defaultContent := map[string]interface{}{
				"type":     "excalidraw",
				"version":  2,
				"source":   "https://excalidraw.com",
				"elements": []interface{}{},
				"appState": map[string]interface{}{
					"viewBackgroundColor": "#ffffff",
					"gridSize":            nil,
				},
			}
			content, _ = json.Marshal(defaultContent)
		}
		metadata.Size = int64(len(content))
	}

	err = h.fileUseCase.CreateFile(metadata, content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error creating file"})
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

	metadata := &models.FileMetadata{
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
