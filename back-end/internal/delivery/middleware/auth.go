package middleware

import (
	"net/http"

	"github.com/gofiber/fiber/v2"

	"myScalidraw/internal/domain/useCase/user"
	"myScalidraw/pkg/jwt"
)

type AuthMiddleware struct {
	jwtManager  *jwt.JWTManager
	userUseCase *user.UserUseCase
}

func NewAuthMiddleware(jwtManager *jwt.JWTManager, userUseCase *user.UserUseCase) *AuthMiddleware {
	return &AuthMiddleware{
		jwtManager:  jwtManager,
		userUseCase: userUseCase,
	}
}

func (m *AuthMiddleware) RequireSystemSetup() fiber.Handler {
	return func(c *fiber.Ctx) error {
		hasUsers, err := m.userUseCase.HasSystemUsers()
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
				"error": "Error checking system setup",
			})
		}

		if !hasUsers {
			return c.Status(http.StatusPreconditionRequired).JSON(fiber.Map{
				"error":        "System setup required",
				"message":      "No users found in system. Please create the first user.",
				"setup_needed": true,
			})
		}

		return c.Next()
	}
}

func (m *AuthMiddleware) RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Cookies("auth_token")
		if token == "" {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Authentication required",
			})
		}

		claims, err := m.jwtManager.ValidateToken(token)
		if err != nil {
			return c.Status(http.StatusUnauthorized).JSON(fiber.Map{
				"error": "Invalid or expired token",
			})
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_type", claims.UserType)

		return c.Next()
	}
}

func (m *AuthMiddleware) RequireOwnerOrAdmin() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "owner" && userType != "admin" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Owner or admin access required",
			})
		}
		return c.Next()
	}
}

func (m *AuthMiddleware) RequireOwner() fiber.Handler {
	return func(c *fiber.Ctx) error {
		userType := c.Locals("user_type")
		if userType != "owner" {
			return c.Status(http.StatusForbidden).JSON(fiber.Map{
				"error": "Owner access required",
			})
		}
		return c.Next()
	}
}

func GetUserID(c *fiber.Ctx) string {
	if userID, ok := c.Locals("user_id").(string); ok {
		return userID
	}
	return ""
}

func GetUserEmail(c *fiber.Ctx) string {
	if email, ok := c.Locals("user_email").(string); ok {
		return email
	}
	return ""
}

func GetUserType(c *fiber.Ctx) string {
	if userType, ok := c.Locals("user_type").(string); ok {
		return userType
	}
	return ""
}
