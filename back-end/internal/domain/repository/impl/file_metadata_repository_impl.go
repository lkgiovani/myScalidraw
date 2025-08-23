package impl

import (
	"myScalidraw/infra/database"
	"myScalidraw/internal/domain/models"
)

type FileMetadataRepositoryImpl struct {
	db *database.DB
}

func NewFileMetadataRepository(db *database.DB) *FileMetadataRepositoryImpl {
	return &FileMetadataRepositoryImpl{
		db: db,
	}
}

func (r *FileMetadataRepositoryImpl) GetAll() (models.FileMetadataList, error) {
	var metadata models.FileMetadataList
	result := r.db.Find(&metadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return metadata, nil
}

func (r *FileMetadataRepositoryImpl) GetByID(id string) (*models.FileMetadata, error) {
	var metadata models.FileMetadata
	result := r.db.First(&metadata, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &metadata, nil
}

func (r *FileMetadataRepositoryImpl) GetByParentID(parentID string) (models.FileMetadataList, error) {
	var metadata models.FileMetadataList
	result := r.db.Find(&metadata, "parent_id = ?", parentID)
	if result.Error != nil {
		return nil, result.Error
	}
	return metadata, nil
}

func (r *FileMetadataRepositoryImpl) Create(metadata *models.FileMetadata) error {
	result := r.db.Create(metadata)
	return result.Error
}

func (r *FileMetadataRepositoryImpl) Update(metadata *models.FileMetadata) error {
	result := r.db.Save(metadata)
	return result.Error
}

func (r *FileMetadataRepositoryImpl) Delete(id string) error {
	result := r.db.Unscoped().Delete(&models.FileMetadata{}, "id = ?", id)
	return result.Error
}
