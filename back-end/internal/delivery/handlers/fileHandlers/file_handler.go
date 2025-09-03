package fileHandlers

import (
	"myScalidraw/internal/delivery/middleware"
	"myScalidraw/internal/domain/useCase/file"
)

type FileHandler struct {
	fileUseCase    *file.FileUseCase
	authMiddleware *middleware.AuthMiddleware
}

func NewFileHandler(fileUseCase *file.FileUseCase, authMiddleware *middleware.AuthMiddleware) *FileHandler {
	return &FileHandler{
		fileUseCase:    fileUseCase,
		authMiddleware: authMiddleware,
	}
}
