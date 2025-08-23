package fileHandlers

import (
	"encoding/json"
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
	IsFolder    bool   `json:"isFolder"`
	ContentType string `json:"contentType"`
}

func (h *FileHandler) CreateFile(c *fiber.Ctx) error {
	var params CreateFileRequest
	if err := c.BodyParser(&params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "erro ao parsear corpo da requisição"})
	}

	fileID, err := uuid.GenerateUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao gerar UUID"})
	}

	fileName := params.Name
	if !params.IsFolder && !strings.HasSuffix(fileName, ".excalidraw") {
		fileName = strings.TrimSuffix(fileName, ".json") + ".excalidraw"
	}

	metadata := &models.FileMetadata{
		ID:           fileID,
		ParentID:     params.ParentID,
		Name:         fileName,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
		IsFolder:     params.IsFolder,
		Size:         0,
		ContentType:  "application/vnd.excalidraw+json",
		LastModified: time.Now(),
	}

	var content []byte
	if !params.IsFolder {
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
		metadata.Size = int64(len(content))
	}

	err = h.fileUseCase.CreateFile(metadata, content)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao criar arquivo"})
	}

	return c.JSON(fiber.Map{
		"message": "arquivo criado com sucesso",
		"id":      fileID,
	})
}

func (h *FileHandler) UploadFile(c *fiber.Ctx) error {
	fileContent := c.Body()
	if len(fileContent) == 0 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "arquivo vazio ou não encontrado"})
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(fileContent, &jsonData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error":   "conteúdo deve ser JSON válido",
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao processar JSON"})
	}

	fileName := c.Get("X-File-Name")
	if fileName == "" {
		fileName = "Untitled-" + time.Now().Format("2006-01-02-1504") + ".excalidraw"
	} else if !strings.HasSuffix(fileName, ".excalidraw") {
		fileName = strings.TrimSuffix(fileName, ".json") + ".excalidraw"
	}

	parentID := c.Query("parentId")
	if parentID == "" {
		parentID = "rascunhos"
	}

	contentType := "application/vnd.excalidraw+json"

	fileID, err := uuid.GenerateUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao gerar UUID"})
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
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "erro ao fazer upload do arquivo"})
	}

	return c.JSON(fiber.Map{
		"message": "arquivo enviado com sucesso",
		"id":      fileID,
		"name":    fileName,
		"size":    len(validatedContent),
	})
}
