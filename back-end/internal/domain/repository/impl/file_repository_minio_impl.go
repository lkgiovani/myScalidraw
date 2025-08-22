package impl

import (
	"encoding/json"
	"fmt"
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

	repo.initializeMetadata()

	repo.loadFileSystem()

	return repo
}

func (r *FileRepositoryMinioImpl) initializeMetadata() {

	metadata, err := r.metadataRepo.GetAll()
	if err != nil || len(metadata) == 0 {

		rascunhos := &models.FileMetadata{
			ID:           "rascunhos",
			Name:         "Rascunhos",
			IsFolder:     true,
			ParentID:     "",
			LastModified: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = r.metadataRepo.Create(rascunhos)
		if err != nil {
			fmt.Printf("Erro ao criar pasta raiz: %v\n", err)
		}

		exemplo := &models.FileMetadata{
			ID:           "exemplo-salve",
			Name:         "Untitled-2025-06-30-1107",
			IsFolder:     false,
			ParentID:     "rascunhos",
			ContentType:  "application/json",
			StoragePath:  "exemplo-salve.json",
			LastModified: time.Now(),
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		err = r.metadataRepo.Create(exemplo)
		if err != nil {
			fmt.Printf("Erro ao criar arquivo de exemplo: %v\n", err)
		}

		content, err := loadLocalExcalidrawFile()
		if err == nil {
			_, uploadErr := r.minioClient.UploadFile("exemplo-salve", []byte(content))
			if uploadErr != nil {
				fmt.Printf("Erro ao fazer upload do arquivo para o MinIO: %v\n", uploadErr)
			}
		}
	}
}

func (r *FileRepositoryMinioImpl) loadFileSystem() {
	metadata, err := r.metadataRepo.GetAll()
	if err != nil {
		fmt.Printf("Erro ao carregar metadados: %v\n", err)

		r.fileSystem = []models.FileItem{
			{
				ID:         "rascunhos",
				Name:       "Rascunhos",
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

	if metadata.IsFolder {
		children, err := r.metadataRepo.GetByParentID(id)
		if err == nil {
			for _, child := range children {
				childItem := child.ToFileItem()
				item.Children = append(item.Children, childItem)
			}
		}
	}

	return &item
}

func (r *FileRepositoryMinioImpl) SaveFile(id string, content string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	metadata, err := r.metadataRepo.GetByID(id)
	if err != nil {
		return fmt.Errorf("arquivo não encontrado: %s", id)
	}

	now := time.Now()
	metadata.LastModified = now
	metadata.UpdatedAt = now

	metadata.Size = int64(len(content))

	err = r.metadataRepo.Update(metadata)
	if err != nil {
		return fmt.Errorf("erro ao atualizar metadados: %w", err)
	}

	_, err = r.minioClient.UploadFile(id, []byte(content))
	if err != nil {
		return fmt.Errorf("erro ao salvar arquivo no MinIO: %w", err)
	}

	return nil
}

func (r *FileRepositoryMinioImpl) GetExcalidrawFile() (string, error) {

	content, err := r.minioClient.GetFile("exemplo-salve")
	if err != nil {

		localContent, localErr := loadLocalExcalidrawFile()
		if localErr != nil {
			return "", fmt.Errorf("erro ao carregar arquivo Excalidraw: %w", err)
		}

		_, uploadErr := r.minioClient.UploadFile("exemplo-salve", []byte(localContent))
		if uploadErr != nil {
			return "", fmt.Errorf("erro ao fazer upload do arquivo para o MinIO: %w", uploadErr)
		}

		return localContent, nil
	}

	return string(content), nil
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
