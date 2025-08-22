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
	}
}

type FileMetadataList []*FileMetadata

func (list FileMetadataList) ToFileSystem() []FileItem {

	itemMap := make(map[string]*FileItem)

	for _, metadata := range list {
		item := metadata.ToFileItem()
		itemMap[item.ID] = &item
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
