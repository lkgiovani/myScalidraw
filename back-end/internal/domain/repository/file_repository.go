package repository

import (
	"myScalidraw/internal/domain/models"
)

type FileRepository interface {
	GetFileSystem() []models.FileItem
	GetFileByID(id string) *models.FileItem
	SaveFile(id string, content string) error
	GetExcalidrawFile() (string, error)
}
