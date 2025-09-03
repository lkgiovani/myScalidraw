package file_repository

import "myScalidraw/internal/domain/models/file_model"

type FileMetadataRepository interface {
	GetAll() (file_model.FileMetadataList, error)

	GetByID(id string) (*file_model.FileMetadata, error)

	GetByParentID(parentID string) (file_model.FileMetadataList, error)

	Create(metadata *file_model.FileMetadata) error

	Update(metadata *file_model.FileMetadata) error

	Delete(id string) error
}
