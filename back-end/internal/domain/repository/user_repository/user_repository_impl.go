package user_repository

import (
	"time"

	"myScalidraw/infra/database"
	"myScalidraw/internal/domain/models/user_model"
)

type userRepositoryImpl struct {
	db *database.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepositoryImpl{
		db: db,
	}
}

func (r *userRepositoryImpl) Create(user *user_model.User) error {
	return r.db.Create(user).Error
}

func (r *userRepositoryImpl) GetByID(id string) (*user_model.User, error) {
	var user user_model.User
	err := r.db.Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetByEmail(email string) (*user_model.User, error) {
	var user user_model.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepositoryImpl) GetAll() (user_model.UserList, error) {
	var users []*user_model.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return user_model.UserList(users), nil
}

func (r *userRepositoryImpl) Update(user *user_model.User) error {
	return r.db.Save(user).Error
}

func (r *userRepositoryImpl) Delete(id string) error {
	return r.db.Where("id = ?", id).Delete(&user_model.User{}).Error
}

func (r *userRepositoryImpl) HasAnyUser() (bool, error) {
	var count int64
	err := r.db.Model(&user_model.User{}).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *userRepositoryImpl) CountUsers() (int64, error) {
	var count int64
	err := r.db.Model(&user_model.User{}).Count(&count).Error
	return count, err
}

func (r *userRepositoryImpl) GetActiveUsers(sessionTimeoutMinutes int64) (user_model.UserList, error) {
	var users []*user_model.User
	cutoff := time.Now().Add(-time.Duration(sessionTimeoutMinutes) * time.Minute)

	err := r.db.Where("type IN (?, ?) OR (type = ? AND last_activity > ?)",
		user_model.UserTypeOwner,
		user_model.UserTypeAdmin,
		user_model.UserTypeGuest,
		cutoff).Find(&users).Error

	if err != nil {
		return nil, err
	}
	return user_model.UserList(users), nil
}

func (r *userRepositoryImpl) UpdateLastActivity(userID string) error {
	return r.db.Model(&user_model.User{}).Where("id = ?", userID).Update("last_activity", time.Now()).Error
}
