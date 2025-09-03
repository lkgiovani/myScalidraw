package file_model

type FileItem struct {
	ID           string      `json:"id"`
	Name         string      `json:"name"`
	IsFolder     bool        `json:"isFolder"`
	Children     []FileItem  `json:"children,omitempty"`
	Data         interface{} `json:"data,omitempty"`
	LastModified int64       `json:"lastModified,omitempty"`
	ParentID     string      `json:"parentId,omitempty"`
	IsExpanded   bool        `json:"isExpanded,omitempty"`
	Path         string      `json:"path,omitempty"`
}
