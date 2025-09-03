package userHandlers

import (
	"net/http"
	"time"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	if req.Email == "" || req.Password == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Email and password are required",
		})
	}

	user, err := h.userUseCase.GetUserByEmail(req.Email)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	if !h.userUseCase.ValidatePassword(user, req.Password) {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid email or password",
		})
	}

	err = h.userUseCase.UpdateUserActivity(user.ID)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating user activity",
		})
	}

	token, err := h.jwtManager.GenerateToken(user.ID, user.Email, string(user.Type))
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error generating token",
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "auth_token",
		Value:    token,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
		Secure:   false,
		SameSite: "Lax",
	})

	return c.JSON(fiber.Map{
		"message": "Login successful",
		"user":    user.ToUserInfo(),
	})
}
