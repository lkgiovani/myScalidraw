package userHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) HasUsers(c *fiber.Ctx) error {
	hasUsers, err := h.userUseCase.HasSystemUsers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error checking system users",
		})
	}

	return c.JSON(fiber.Map{
		"hasUsers": hasUsers,
	})
}
