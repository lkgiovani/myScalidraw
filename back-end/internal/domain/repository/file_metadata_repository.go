package repository

import (
	"myScalidraw/internal/domain/models"
)

type FileMetadataRepository interface {
	GetAll() (models.FileMetadataList, error)

	GetByID(id string) (*models.FileMetadata, error)

	GetByParentID(parentID string) (models.FileMetadataList, error)

	Create(metadata *models.FileMetadata) error

	Update(metadata *models.FileMetadata) error

	Delete(id string) error
}
