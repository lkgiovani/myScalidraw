package repository

import (
	"myScalidraw/internal/domain/models"
)

type FileRepository interface {
	GetFileSystem() []models.FileItem
	GetFileByID(id string) *models.FileItem
	SaveFile(id string, content string) error
	GetFileContent(id string) (string, error)
	UploadFile(id string, content []byte) error
	DeleteFile(id string) error
	RenameFile(id string, newName string) error
}
