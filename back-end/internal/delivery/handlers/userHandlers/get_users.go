package userHandlers

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

func (h *UserHandler) GetUsers(c *fiber.Ctx) error {
	users, err := h.userUseCase.GetAllUsers()
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching users",
		})
	}

	return c.JSON(fiber.Map{
		"users": users.ToUserInfoList(),
	})
}
