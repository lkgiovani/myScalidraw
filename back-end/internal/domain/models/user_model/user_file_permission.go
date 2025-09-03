package user_model

import (
	"time"

	"gorm.io/gorm"
)

type PermissionType string

const (
	PermissionTypeRead  PermissionType = "read"
	PermissionTypeWrite PermissionType = "write"
	PermissionTypeAdmin PermissionType = "admin"
)

type UserFilePermission struct {
	ID         string         `json:"id" gorm:"primaryKey"`
	UserID     string         `json:"userId" gorm:"index"`
	FileID     string         `json:"fileId" gorm:"index"`
	Permission PermissionType `json:"permission"`
	GrantedBy  string         `json:"grantedBy"`
	ExpiresAt  *time.Time     `json:"expiresAt,omitempty"`
	CreatedAt  time.Time      `json:"createdAt"`
	UpdatedAt  time.Time      `json:"updatedAt"`
	DeletedAt  gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (ufp *UserFilePermission) IsExpired() bool {
	if ufp.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*ufp.ExpiresAt)
}

func (ufp *UserFilePermission) CanRead() bool {
	return !ufp.IsExpired() && (ufp.Permission == PermissionTypeRead || ufp.Permission == PermissionTypeWrite || ufp.Permission == PermissionTypeAdmin)
}

func (ufp *UserFilePermission) CanWrite() bool {
	return !ufp.IsExpired() && (ufp.Permission == PermissionTypeWrite || ufp.Permission == PermissionTypeAdmin)
}

func (ufp *UserFilePermission) CanAdmin() bool {
	return !ufp.IsExpired() && ufp.Permission == PermissionTypeAdmin
}

func (ufp *UserFilePermission) ToPermissionInfo() UserFilePermissionInfo {
	var expiresAt *int64
	if ufp.ExpiresAt != nil {
		timestamp := ufp.ExpiresAt.Unix() * 1000
		expiresAt = &timestamp
	}

	return UserFilePermissionInfo{
		ID:         ufp.ID,
		UserID:     ufp.UserID,
		FileID:     ufp.FileID,
		Permission: ufp.Permission,
		GrantedBy:  ufp.GrantedBy,
		ExpiresAt:  expiresAt,
		CreatedAt:  ufp.CreatedAt.Unix() * 1000,
	}
}

type UserFilePermissionInfo struct {
	ID         string         `json:"id"`
	UserID     string         `json:"userId"`
	FileID     string         `json:"fileId"`
	Permission PermissionType `json:"permission"`
	GrantedBy  string         `json:"grantedBy"`
	ExpiresAt  *int64         `json:"expiresAt,omitempty"`
	CreatedAt  int64          `json:"createdAt"`
}

type UserFilePermissionList []*UserFilePermission

func (list UserFilePermissionList) ToPermissionInfoList() []UserFilePermissionInfo {
	var permissions []UserFilePermissionInfo

	for _, permission := range list {
		if permission == nil {
			continue
		}

		permissions = append(permissions, permission.ToPermissionInfo())
	}

	return permissions
}

func (list UserFilePermissionList) FilterByUser(userID string) UserFilePermissionList {
	var filtered UserFilePermissionList

	for _, permission := range list {
		if permission != nil && permission.UserID == userID {
			filtered = append(filtered, permission)
		}
	}

	return filtered
}

func (list UserFilePermissionList) FilterByFile(fileID string) UserFilePermissionList {
	var filtered UserFilePermissionList

	for _, permission := range list {
		if permission != nil && permission.FileID == fileID {
			filtered = append(filtered, permission)
		}
	}

	return filtered
}

func (list UserFilePermissionList) FilterActive() UserFilePermissionList {
	var active UserFilePermissionList

	for _, permission := range list {
		if permission != nil && !permission.IsExpired() {
			active = append(active, permission)
		}
	}

	return active
}
