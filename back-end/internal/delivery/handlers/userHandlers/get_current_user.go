package userHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetCurrentUser(c *fiber.Ctx) error {
	token := c.Cookies("auth_token")
	if token == "" {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "No authentication token",
		})
	}

	claims, err := h.jwtManager.ValidateToken(token)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired token",
		})
	}

	user, err := h.userUseCase.GetUserByID(claims.UserID)
	if err != nil {
		return c.Status(http.StatusNotFound).JSON(fiber.Map{
			"error": "User not found",
		})
	}

	return c.JSON(fiber.Map{
		"user": user.ToUserInfo(),
	})
}
