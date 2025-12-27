package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/secure-notes/internal/http/response"
	"github.com/secure-notes/internal/security"
	"strings"
)

const LocalUserIDKey = "user_id"

func AuthRequired(jwtm *security.JWTManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		h := c.Get("Authorization")
		if h == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(response.NewError(response.CodeUnauthorized, "missing authorization header"))
		}

		parts := strings.SplitN(h, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || parts[1] == "" {
			return c.Status(fiber.StatusUnauthorized).
				JSON(response.NewError(response.CodeUnauthorized, "invalid authorization header"))
		}

		userID, err := jwtm.Parse(parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).
				JSON(response.NewError(response.CodeUnauthorized, "invalid token"))
		}

		c.Locals(LocalUserIDKey, userID)
		return c.Next()
	}
}
