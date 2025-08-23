package impl

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"myScalidraw/infra/database"
	"myScalidraw/infra/storage"
	"myScalidraw/internal/domain/models"
	"myScalidraw/internal/domain/repository"
)

type FileRepositoryMinioImpl struct {
	mu           sync.Mutex
	fileSystem   []models.FileItem
	minioClient  *storage.MinIO
	metadataRepo repository.FileMetadataRepository
	db           *database.DB
}

func NewFileRepositoryMinio(minioClient *storage.MinIO, db *database.DB, metadataRepo repository.FileMetadataRepository) *FileRepositoryMinioImpl {
	repo := &FileRepositoryMinioImpl{
		minioClient:  minioClient,
		metadataRepo: metadataRepo,
		db:           db,
	}

	repo.loadFileSystem()

	return repo
}

func (r *FileRepositoryMinioImpl) loadFileSystem() {
	metadata, err := r.metadataRepo.GetAll()
	if err != nil {
		fmt.Printf("Error loading metadata: %v\n", err)

		r.fileSystem = []models.FileItem{
			{
				ID:         "drafts",
				Name:       "Drafts",
				IsFolder:   true,
				IsExpanded: true,
				Children:   []models.FileItem{},
			},
		}
		return
	}

	r.fileSystem = metadata.ToFileSystem()
}

func (r *FileRepositoryMinioImpl) GetFileSystem() []models.FileItem {

	r.loadFileSystem()
	return r.fileSystem
}

func (r *FileRepositoryMinioImpl) findFileByID(items []models.FileItem, id string) *models.FileItem {
	for i := range items {
		if items[i].ID == id {
			return &items[i]
		}
		if items[i].IsFolder && len(items[i].Children) > 0 {
			if found := r.findFileByID(items[i].Children, id); found != nil {
				return found
			}
		}
	}
	return nil
}

func (r *FileRepositoryMinioImpl) GetFileByID(id string) *models.FileItem {

	metadata, err := r.metadataRepo.GetByID(id)
	if err != nil {

		return r.findFileByID(r.fileSystem, id)
	}

	item := metadata.ToFileItem()

	item.Path = r.buildItemPath(metadata)

	if metadata.IsFolder {
		children, err := r.metadataRepo.GetByParentID(id)
		if err == nil {
			for _, child := range children {
				childItem := child.ToFileItem()
				childItem.Path = r.buildItemPath(child)
				item.Children = append(item.Children, childItem)
			}
		}
	}

	return &item
}

func (r *FileRepositoryMinioImpl) buildItemPath(metadata *models.FileMetadata) string {
	var pathParts []string
	current := metadata

	for current != nil && current.ParentID != "" {
		parent, err := r.metadataRepo.GetByID(current.ParentID)
		if err == nil {
			pathParts = append([]string{parent.Name}, pathParts...)
			current = parent
		} else {
			break
		}
	}

	if len(pathParts) == 0 {
		return "/"
	}

	result := "/"
	for i, part := range pathParts {
		if i > 0 {
			result += "/"
		}
		result += part
	}
	return result
}

func (r *FileRepositoryMinioImpl) SaveFile(id string, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata, err := r.metadataRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %s", id)
	}

	now := time.Now()
	metadata.LastModified = now
	metadata.UpdatedAt = now

	metadata.Size = int64(len(content))

	err = r.metadataRepo.Update(metadata)
	if err != nil {
		return fmt.Errorf("error updating metadata: %w", err)
	}

	_, err = r.minioClient.UploadFile(id, []byte(content))
	if err != nil {
		return fmt.Errorf("error saving file to MinIO: %w", err)
	}

	return nil
}

func (r *FileRepositoryMinioImpl) GetFileContent(id string) (string, error) {
	content, err := r.minioClient.GetFile(id)
	if err != nil {
		if id == "exemplo-salve" {
			localContent, localErr := loadLocalExcalidrawFile()
			if localErr != nil {
				return "", fmt.Errorf("error loading file: %w", err)
			}

			_, uploadErr := r.minioClient.UploadFile(id, []byte(localContent))
			if uploadErr != nil {
				return "", fmt.Errorf("error uploading file to MinIO: %w", uploadErr)
			}

			return localContent, nil
		}
		return "", fmt.Errorf("error fetching file: %w", err)
	}

	return string(content), nil
}

func (r *FileRepositoryMinioImpl) UploadFile(id string, content []byte) error {
	_, err := r.minioClient.UploadFile(id, content)
	if err != nil {
		return fmt.Errorf("error uploading file: %w", err)
	}
	return nil
}

func (r *FileRepositoryMinioImpl) DeleteFile(id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata, err := r.metadataRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %s", id)
	}

	if metadata.IsFolder {
		children, childErr := r.metadataRepo.GetByParentID(id)
		if childErr == nil && len(children) > 0 {
			for _, child := range children {
				if deleteErr := r.DeleteFile(child.ID); deleteErr != nil {
					return fmt.Errorf("error deleting child file %s: %w", child.ID, deleteErr)
				}
			}
		}
	} else {
		if minioErr := r.minioClient.DeleteFile(id); minioErr != nil {
			return fmt.Errorf("error deleting file from MinIO: %w", minioErr)
		}
	}

	err = r.metadataRepo.Delete(id)
	if err != nil {
		return fmt.Errorf("error deleting metadata: %w", err)
	}

	return nil
}

func (r *FileRepositoryMinioImpl) RenameFile(id string, newName string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata, err := r.metadataRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("file not found: %s", id)
	}

	if !metadata.IsFolder && !strings.HasSuffix(newName, ".excalidraw") {
		newName = strings.TrimSuffix(newName, ".json") + ".excalidraw"
	}

	metadata.Name = newName
	metadata.UpdatedAt = time.Now()

	err = r.metadataRepo.Update(metadata)
	if err != nil {
		return fmt.Errorf("error updating metadata: %w", err)
	}

	return nil
}

func loadLocalExcalidrawFile() (string, error) {

	content, err := loadFileFromDisk("./Untitled-2025-06-30-1107.excalidraw")
	if err != nil {

		emptyExcalidraw := map[string]interface{}{
			"type":     "excalidraw",
			"version":  2,
			"source":   "https://excalidraw.com",
			"elements": []interface{}{},
			"appState": map[string]interface{}{
				"viewBackgroundColor": "#ffffff",
				"gridSize":            nil,
			},
		}

		jsonContent, err := json.Marshal(emptyExcalidraw)
		if err != nil {
			return "", err
		}

		return string(jsonContent), nil
	}

	return content, nil
}

func loadFileFromDisk(path string) (string, error) {

	content, err := storage.ReadFileFromDisk(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
