package file_repository

import (
	"myScalidraw/infra/database"
	"myScalidraw/internal/domain/models/file_model"
)

type FileMetadataRepositoryImpl struct {
	db *database.DB
}

func NewFileMetadataRepository(db *database.DB) FileMetadataRepository {
	return &FileMetadataRepositoryImpl{
		db: db,
	}
}

func (r *FileMetadataRepositoryImpl) GetAll() (file_model.FileMetadataList, error) {
	var metadata file_model.FileMetadataList
	result := r.db.Find(&metadata)
	if result.Error != nil {
		return nil, result.Error
	}
	return metadata, nil
}

func (r *FileMetadataRepositoryImpl) GetByID(id string) (*file_model.FileMetadata, error) {
	var metadata file_model.FileMetadata
	result := r.db.First(&metadata, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &metadata, nil
}

func (r *FileMetadataRepositoryImpl) GetByParentID(parentID string) (file_model.FileMetadataList, error) {
	var metadata file_model.FileMetadataList
	result := r.db.Find(&metadata, "parent_id = ?", parentID)
	if result.Error != nil {
		return nil, result.Error
	}
	return metadata, nil
}

func (r *FileMetadataRepositoryImpl) Create(metadata *file_model.FileMetadata) error {
	result := r.db.Create(metadata)
	return result.Error
}

func (r *FileMetadataRepositoryImpl) Update(metadata *file_model.FileMetadata) error {
	result := r.db.Save(metadata)
	return result.Error
}

func (r *FileMetadataRepositoryImpl) Delete(id string) error {
	result := r.db.Unscoped().Delete(&file_model.FileMetadata{}, "id = ?", id)
	return result.Error
}
