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

type CreateFileRequest struct {
	Name        string `json:"name"`
	ParentID    string `json:"parentId"`
	ParentPath  string `json:"parentPath"`
	IsFolder    bool   `json:"isFolder"`
	ContentType string `json:"contentType"`
	Content     string `json:"content"`
}

func (h *FileHandler) CreateFile(c *fiber.Ctx) error {

	var params CreateFileRequest
	if err := c.BodyParser(&params); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "error parsing request body"})
	}

	fileID, err := uuid.GenerateUUID()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "error generating UUID"})
	}

	fileName := params.Name

	isFolder := params.IsFolder
	if params.Content != "" {
		isFolder = false
	}

	if !isFolder && !strings.HasSuffix(fileName, ".excalidraw") {
		fileName = strings.TrimSuffix(fileName, ".json") + ".excalidraw"
	}

	var storagePath string
	if params.ParentID != "" {
		parent, parentErr := h.fileUseCase.GetFileByID(params.ParentID)
		if parentErr == nil && parent != nil {

			parentPath := strings.TrimSuffix(parent.Path, "/")
			if parentPath == "" {
				parentPath = "/"
			}
			storagePath = parentPath + "/" + fileName
		} else {

			storagePath = "/" + fileName
		}
	} else {
		storagePath = "/" + fileName
	}

	metadata := &file_model.FileMetadata{
		ID:           fileID,
		ParentID:     params.ParentID,
		Name:         fileName,
		StoragePath:  storagePath,
		Path:         storagePath,
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

	fileItem := metadata.ToFileItem()

	if metadata.IsFolder {
		fileItem.Children = []file_model.FileItem{}
	}

	return c.JSON(fileItem)
}
