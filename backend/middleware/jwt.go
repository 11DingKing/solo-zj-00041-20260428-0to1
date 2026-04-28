package middleware

import (
	"strings"

	"canteen-system/models"
	"canteen-system/utils"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(401).JSON(models.Error(401, "未授权访问"))
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			return c.Status(401).JSON(models.Error(401, "认证格式错误"))
		}

		claims, err := utils.ParseToken(parts[1])
		if err != nil {
			return c.Status(401).JSON(models.Error(401, "令牌无效或已过期"))
		}

		c.Locals("user_id", claims.UserID)
		c.Locals("username", claims.Username)
		c.Locals("role", claims.Role)

		return c.Next()
	}
}

func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("role").(string)
		if !ok {
			return c.Status(403).JSON(models.Error(403, "权限不足"))
		}

		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return c.Status(403).JSON(models.Error(403, "权限不足"))
	}
}
