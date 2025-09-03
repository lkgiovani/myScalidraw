package file_repository

import (
	"myScalidraw/internal/domain/models/file_model"
)

type FileRepository interface {
	GetFileSystem() []file_model.FileItem
	GetFileByID(id string) *file_model.FileItem
	SaveFile(id string, content string) error
	GetFileContent(id string) (string, error)
	UploadFile(id string, content []byte) error
	CreateFolder(folderPath string) error
	DeleteFile(id string) error
	RenameFile(id string, newName string) error
}
