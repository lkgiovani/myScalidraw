package user_repository

import "myScalidraw/internal/domain/models/user_model"

type UserRepository interface {
	Create(user *user_model.User) error
	GetByID(id string) (*user_model.User, error)
	GetByEmail(email string) (*user_model.User, error)
	GetAll() (user_model.UserList, error)
	Update(user *user_model.User) error
	Delete(id string) error
	HasAnyUser() (bool, error)
	CountUsers() (int64, error)
	GetActiveUsers(sessionTimeout int64) (user_model.UserList, error)
	UpdateLastActivity(userID string) error
}
