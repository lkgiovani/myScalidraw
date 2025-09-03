package file

import (
	"encoding/json"
	"myScalidraw/internal/domain/models/file_model"
	"myScalidraw/internal/domain/repository/file_repository"
)

type FileUseCase struct {
	fileRepo     file_repository.FileRepository
	metadataRepo file_repository.FileMetadataRepository
}

func NewFileUseCase(fileRepo file_repository.FileRepository, metadataRepo file_repository.FileMetadataRepository) *FileUseCase {
	return &FileUseCase{
		fileRepo:     fileRepo,
		metadataRepo: metadataRepo,
	}
}

func (uc *FileUseCase) GetFiles() []file_model.FileItem {

	metadata, err := uc.metadataRepo.GetAll()
	if err != nil {

		return []file_model.FileItem{}
	}

	flatList := metadata.ToFlatList()

	return flatList
}

func (uc *FileUseCase) GetFileByID(id string) (*file_model.FileItem, error) {
	file := uc.fileRepo.GetFileByID(id)
	if file == nil {
		return nil, nil
	}

	if !file.IsFolder {
		content, err := uc.fileRepo.GetFileContent(id)
		if err == nil {
			var data map[string]interface{}
			if err := json.Unmarshal([]byte(content), &data); err == nil {
				file.Data = data
			}
		}
	}

	return file, nil
}

func (uc *FileUseCase) SaveFile(id string, content string) error {
	return uc.fileRepo.SaveFile(id, content)
}

func (uc *FileUseCase) CreateFile(metadata *file_model.FileMetadata, content []byte) error {
	err := uc.metadataRepo.Create(metadata)
	if err != nil {
		return err
	}

	if metadata.IsFolder {
		return uc.fileRepo.CreateFolder(metadata.Path)
	}

	if len(content) > 0 {
		return uc.fileRepo.UploadFile(metadata.ID, content)
	}

	return nil
}

func (uc *FileUseCase) DeleteFile(id string) error {
	return uc.fileRepo.DeleteFile(id)
}

func (uc *FileUseCase) RenameFile(id string, newName string) error {
	return uc.fileRepo.RenameFile(id, newName)
}

func (uc *FileUseCase) GetFileContent(id string) (string, error) {
	return uc.fileRepo.GetFileContent(id)
}
