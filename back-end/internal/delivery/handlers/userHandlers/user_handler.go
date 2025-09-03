package userHandlers

import (
	"myScalidraw/internal/domain/useCase/user"
	"myScalidraw/pkg/jwt"
)

type UserHandler struct {
	userUseCase *user.UserUseCase
	jwtManager  *jwt.JWTManager
}

func NewUserHandler(userUseCase *user.UserUseCase, jwtManager *jwt.JWTManager) *UserHandler {
	return &UserHandler{
		userUseCase: userUseCase,
		jwtManager:  jwtManager,
	}
}

type CreateFirstUserRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}
