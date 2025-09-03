package userHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) CreateFirstUser(c *fiber.Ctx) error {
	var req CreateFirstUserRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" || req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Name, email and password are required",
		})
	}

	if len(req.Password) < 8 {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Password must be at least 8 characters long",
		})
	}

	hasUsers, err := h.userUseCase.HasSystemUsers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking system users",
		})
	}

	if hasUsers {
		return c.Status(http.StatusConflict).JSON(fiber.Map{
			"error": "System already has users. First user was already created.",
		})
	}

	user, err := h.userUseCase.CreateFirstUser(req.Name, req.Email, req.Password)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error creating first user: " + err.Error(),
		})
	}

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"message": "First user created successfully",
		"user":    user.ToUserInfo(),
	})
}
