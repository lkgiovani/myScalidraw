package models

type FileItem struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	IsFolder     bool        `json:"isFolder,omitempty"`
	Children     []FileItem  `json:"children,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	LastModified int64       `json:"lastModified,omitempty"`
	ParentID     string      `json:"parentId,omitempty"`
	IsExpanded   bool        `json:"isExpanded,omitempty"`
}
