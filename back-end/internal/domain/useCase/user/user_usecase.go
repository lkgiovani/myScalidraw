package user

import (
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"

	"myScalidraw/internal/domain/models/user_model"
	"myScalidraw/internal/domain/repository/user_repository"

	"myScalidraw/pkg/uuid"
)

type UserUseCase struct {
	userRepo user_repository.UserRepository
}

func NewUserUseCase(userRepo user_repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		userRepo: userRepo,
	}
}

func (uc *UserUseCase) CreateFirstUser(name, email, password string) (*user_model.User, error) {
	hasUsers, err := uc.userRepo.HasAnyUser()
	if err != nil {
		return nil, err
	}

	if hasUsers {
		return nil, errors.New("system already has users - first user already created")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	userID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	user := &user_model.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		Password:     string(hashedPassword),
		Type:         user_model.UserTypeOwner,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) CreateGuestUser(name, email string, guestType user_model.GuestType, sessionID *string) (*user_model.User, error) {
	userID, err := uuid.GenerateUUID()
	if err != nil {
		return nil, err
	}

	user := &user_model.User{
		ID:           userID,
		Name:         name,
		Email:        email,
		Type:         user_model.UserTypeGuest,
		GuestType:    &guestType,
		SessionID:    sessionID,
		LastActivity: time.Now(),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	err = uc.userRepo.Create(user)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (uc *UserUseCase) GetUserByID(id string) (*user_model.User, error) {
	return uc.userRepo.GetByID(id)
}

func (uc *UserUseCase) GetUserByEmail(email string) (*user_model.User, error) {
	return uc.userRepo.GetByEmail(email)
}

func (uc *UserUseCase) ValidatePassword(user *user_model.User, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

func (uc *UserUseCase) UpdateUserActivity(userID string) error {
	return uc.userRepo.UpdateLastActivity(userID)
}

func (uc *UserUseCase) GetActiveUsers(sessionTimeoutMinutes int64) (user_model.UserList, error) {
	return uc.userRepo.GetActiveUsers(sessionTimeoutMinutes)
}

func (uc *UserUseCase) HasSystemUsers() (bool, error) {
	return uc.userRepo.HasAnyUser()
}

func (uc *UserUseCase) GetAllUsers() (user_model.UserList, error) {
	return uc.userRepo.GetAll()
}

func (uc *UserUseCase) DeleteUser(id string) error {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return err
	}

	if user.IsOwner() {
		count, err := uc.userRepo.CountUsers()
		if err != nil {
			return err
		}
		if count <= 1 {
			return errors.New("cannot delete the only owner user")
		}
	}

	return uc.userRepo.Delete(id)
}
