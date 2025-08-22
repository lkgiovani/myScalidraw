package file

import (
	"encoding/json"

	"myScalidraw/internal/domain/models"
	"myScalidraw/internal/domain/repository"
)

type FileUseCase struct {
	repo repository.FileRepository
}

func NewFileUseCase(repo repository.FileRepository) *FileUseCase {
	return &FileUseCase{
		repo: repo,
	}
}

func (uc *FileUseCase) GetFiles() []models.FileItem {
	return uc.repo.GetFileSystem()
}

func (uc *FileUseCase) GetFileByID(id string) (*models.FileItem, error) {
	file := uc.repo.GetFileByID(id)
	if file == nil {
		return nil, nil
	}

	if !file.IsFolder {

		excalidrawData := map[string]interface{}{
			"elements": []interface{}{},
			"appState": map[string]interface{}{
				"viewBackgroundColor":        "#ffffff",
				"currentItemStrokeColor":     "#000000",
				"currentItemBackgroundColor": "#ffffff",
				"currentItemFillStyle":       "solid",
				"currentItemStrokeWidth":     1,
				"currentItemStrokeStyle":     "solid",
				"currentItemRoughness":       1,
				"currentItemOpacity":         100,
				"currentItemFontFamily":      1,
				"currentItemFontSize":        20,
				"currentItemTextAlign":       "left",
				"currentItemStartArrowhead":  nil,
				"currentItemEndArrowhead":    nil,
				"scrollX":                    0,
				"scrollY":                    0,
				"zoom":                       map[string]interface{}{"value": 1},
				"currentItemRoundness":       "round",
				"gridSize":                   nil,
				"colorPalette":               map[string]interface{}{},
			},
		}

		if id == "exemplo-salve" {
			content, err := uc.repo.GetExcalidrawFile()
			if err == nil {
				var data map[string]interface{}
				if err := json.Unmarshal([]byte(content), &data); err == nil {
					excalidrawData = data
				}
			}
		}

		file.Data = excalidrawData
	}

	return file, nil
}

func (uc *FileUseCase) GetSala(id string) (string, bool) {
	file := uc.repo.GetFileByID(id)
	if file == nil {
		return "", false
	}

	content, err := uc.repo.GetExcalidrawFile()
	if err != nil {
		return "", false
	}

	return content, true
}

func (uc *FileUseCase) SaveSala(id string, content string) error {
	return uc.repo.SaveFile(id, content)
}
