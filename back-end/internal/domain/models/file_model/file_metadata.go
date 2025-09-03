package file_model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

type FileMetadata struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name"`
	IsFolder     bool           `json:"isFolder"`
	ParentID     string         `json:"parentId"`
	StoragePath  string         `json:"storagePath"`
	Path         string         `json:"path"`
	ContentType  string         `json:"contentType"`
	Size         int64          `json:"size"`
	LastModified time.Time      `json:"lastModified"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (fm *FileMetadata) ToFileItem() FileItem {
	return FileItem{
		ID:           fm.ID,
		Name:         fm.Name,
		IsFolder:     fm.IsFolder,
		ParentID:     fm.ParentID,
		LastModified: fm.LastModified.Unix() * 1000,
		Path:         fm.Path,
	}
}

type FileMetadataList []*FileMetadata

func (list FileMetadataList) ToFileSystem() []FileItem {

	itemMap := make(map[string]*FileItem)
	metadataMap := make(map[string]*FileMetadata)

	for _, metadata := range list {
		if metadata == nil {

			continue
		}

		item := metadata.ToFileItem()

		if strings.Contains(item.Path, "//") {
			item.Path = strings.ReplaceAll(item.Path, "//", "/")

		}

		itemMap[item.ID] = &item
		metadataMap[metadata.ID] = metadata

	}

	var rootItems []FileItem

	for _, item := range itemMap {
		if item.ParentID == "" || item.ParentID == "null" {

			rootItems = append(rootItems, *item)
		} else {
			if parent, ok := itemMap[item.ParentID]; ok {
				if parent.Children == nil {
					parent.Children = []FileItem{}
				}
				parent.Children = append(parent.Children, *item)

			} else {

				rootItems = append(rootItems, *item)
			}
		}
	}

	return rootItems
}

func (list FileMetadataList) ToFlatList() []FileItem {

	var items []FileItem

	for _, metadata := range list {
		if metadata == nil {

			continue
		}

		item := metadata.ToFileItem()

		if strings.Contains(item.Path, "//") {
			item.Path = strings.ReplaceAll(item.Path, "//", "/")

		}

		items = append(items, item)

	}

	return items
}
