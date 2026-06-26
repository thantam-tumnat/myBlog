package middlewares

import (
	"strings"

	"blogservice/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

// Locals keys shared with downstream handlers.
const (
	UserIDKey = "userId"
	RoleKey   = "role"
)

// JWTAuth verifies the Bearer token and stores user_id/role in the request
// context. Handlers read c.Locals(UserIDKey) instead of trusting URL params.
func JWTAuth(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "missing authorization header",
			})
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid authorization header format",
			})
		}

		claims, err := utils.ParseAccessToken(secret, parts[1])
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"message": "invalid or expired token",
			})
		}

		c.Locals(UserIDKey, claims.UserID)
		c.Locals(RoleKey, claims.Role)
		return c.Next()
	}
}

// RequireRole guards a route so only users with the given role may proceed.
// Must be chained after JWTAuth.
func RequireRole(role string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Locals(RoleKey) != role {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"message": "forbidden: insufficient permission",
			})
		}
		return c.Next()
	}
}
