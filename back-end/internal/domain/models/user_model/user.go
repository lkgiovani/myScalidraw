package user_model

import (
	"time"

	"gorm.io/gorm"
)

type UserType string
type GuestType string

const (
	UserTypeOwner UserType = "owner"
	UserTypeAdmin UserType = "admin"
	UserTypeGuest UserType = "guest"
)

const (
	GuestTypePersistent GuestType = "persistent"
	GuestTypeTemporary  GuestType = "temporary"
)

type User struct {
	ID           string         `json:"id" gorm:"primaryKey"`
	Name         string         `json:"name"`
	Email        string         `json:"email" gorm:"unique"`
	Password     string         `json:"-" gorm:"not null"`
	Type         UserType       `json:"type"`
	GuestType    *GuestType     `json:"guestType,omitempty"`
	SessionID    *string        `json:"sessionId,omitempty"`
	LastActivity time.Time      `json:"lastActivity"`
	CreatedAt    time.Time      `json:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt"`
	DeletedAt    gorm.DeletedAt `json:"deletedAt" gorm:"index"`
}

func (u *User) IsOwner() bool {
	return u.Type == UserTypeOwner
}

func (u *User) IsAdmin() bool {
	return u.Type == UserTypeAdmin
}

func (u *User) IsGuest() bool {
	return u.Type == UserTypeGuest
}

func (u *User) IsPersistentGuest() bool {
	return u.Type == UserTypeGuest && u.GuestType != nil && *u.GuestType == GuestTypePersistent
}

func (u *User) IsTemporaryGuest() bool {
	return u.Type == UserTypeGuest && u.GuestType != nil && *u.GuestType == GuestTypeTemporary
}

func (u *User) CanAlwaysAccess() bool {
	return u.IsOwner() || u.IsAdmin() || u.IsPersistentGuest()
}

func (u *User) RequiresActiveSession() bool {
	return u.IsTemporaryGuest()
}

func (u *User) UpdateActivity() {
	u.LastActivity = time.Now()
}

func (u *User) ToUserInfo() UserInfo {
	return UserInfo{
		ID:           u.ID,
		Name:         u.Name,
		Email:        u.Email,
		Type:         u.Type,
		GuestType:    u.GuestType,
		LastActivity: u.LastActivity.Unix() * 1000,
	}
}

type UserInfo struct {
	ID           string     `json:"id"`
	Name         string     `json:"name"`
	Email        string     `json:"email"`
	Type         UserType   `json:"type"`
	GuestType    *GuestType `json:"guestType,omitempty"`
	LastActivity int64      `json:"lastActivity"`
}

type UserList []*User

func (list UserList) ToUserInfoList() []UserInfo {
	var userInfos []UserInfo

	for _, user := range list {
		if user == nil {
			continue
		}

		userInfos = append(userInfos, user.ToUserInfo())
	}

	return userInfos
}

func (list UserList) FilterByType(userType UserType) UserList {
	var filtered UserList

	for _, user := range list {
		if user != nil && user.Type == userType {
			filtered = append(filtered, user)
		}
	}

	return filtered
}

func (list UserList) FilterActiveUsers(sessionTimeout time.Duration) UserList {
	var active UserList
	cutoff := time.Now().Add(-sessionTimeout)

	for _, user := range list {
		if user == nil {
			continue
		}

		if user.CanAlwaysAccess() || user.LastActivity.After(cutoff) {
			active = append(active, user)
		}
	}

	return active
}
