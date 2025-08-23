package models

import (
	"time"

	"gorm.io/gorm"
)

type FileMetadata struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name"`
	IsFolder     bool           `json:"isFolder"`
	ParentID     string         `json:"parentId"`
	StoragePath  string         `json:"storagePath"`
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
		Path:         fm.StoragePath,
	}
}

type FileMetadataList []*FileMetadata

func (list FileMetadataList) ToFileSystem() []FileItem {

	itemMap := make(map[string]*FileItem)
	metadataMap := make(map[string]*FileMetadata)

	for _, metadata := range list {
		item := metadata.ToFileItem()
		itemMap[item.ID] = &item
		metadataMap[metadata.ID] = metadata
	}

	buildPath := func(itemID string) string {
		var pathParts []string
		current := metadataMap[itemID]

		for current != nil && current.ParentID != "" {
			if parent, exists := metadataMap[current.ParentID]; exists {
				pathParts = append([]string{parent.Name}, pathParts...)
				current = parent
			} else {
				break
			}
		}

		if len(pathParts) == 0 {
			return "/"
		}
		return "/" + joinPath(pathParts)
	}

	for id, item := range itemMap {
		item.Path = buildPath(id)
	}

	var rootItems []FileItem

	for id, item := range itemMap {
		if item.ParentID == "" {

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

		delete(itemMap, id)
	}

	return rootItems
}

func joinPath(parts []string) string {
	if len(parts) == 0 {
		return ""
	}
	if len(parts) == 1 {
		return parts[0]
	}

	result := parts[0]
	for i := 1; i < len(parts); i++ {
		result += "/" + parts[i]
	}
	return result
}
