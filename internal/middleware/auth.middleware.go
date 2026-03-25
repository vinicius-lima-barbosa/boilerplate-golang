package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/responses"
	"github.com/vinicius-lima-barbosa/boilerplate-golang/internal/services"
)

const AuthenticatedUserIDKey = "auth_user_id"
const AuthenticatedUserTypeKey = "auth_user_type"

func JWTAuthMiddleware(authService services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" {
			return responses.NewErrorResponse(c, fiber.StatusUnauthorized, "Missing authorization header", nil)
		}

		parts := strings.SplitN(header, " ", 2)
		if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
			return responses.NewErrorResponse(c, fiber.StatusUnauthorized, "Invalid authorization header format", nil)
		}

		claims, err := authService.ValidateAccessToken(parts[1])
		if err != nil {
			return responses.NewErrorResponse(c, fiber.StatusUnauthorized, "Invalid or expired access token", nil)
		}

		c.Locals(AuthenticatedUserIDKey, claims.UserID)
		c.Locals(AuthenticatedUserTypeKey, claims.UserType)

		return c.Next()
	}
}
