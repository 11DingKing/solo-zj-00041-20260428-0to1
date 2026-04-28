package handlers

import (
	"canteen-system/database"
	"canteen-system/models"
	"canteen-system/utils"

	"github.com/gofiber/fiber/v2"
)

func Login(c *fiber.Ctx) error {
	var req models.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	var user models.User
	if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
		return c.Status(401).JSON(models.Error(401, "用户名或密码错误"))
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return c.Status(401).JSON(models.Error(401, "用户名或密码错误"))
	}

	token, err := utils.GenerateToken(&user)
	if err != nil {
		return c.Status(500).JSON(models.Error(500, "生成令牌失败"))
	}

	var allergens []models.UserAllergen
	database.DB.Where("user_id = ?", user.ID).Find(&allergens)
	user.Allergens = allergens

	return c.JSON(models.Success(models.LoginResponse{
		Token: token,
		User:  &user,
	}))
}

func GetCurrentUser(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "用户不存在"))
	}

	var allergens []models.UserAllergen
	database.DB.Where("user_id = ?", user.ID).Find(&allergens)
	user.Allergens = allergens

	return c.JSON(models.Success(user))
}

func UpdateUserAllergens(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req models.UpdateAllergensRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if err := database.DB.Where("user_id = ?", userID).Delete(&models.UserAllergen{}).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "更新过敏原失败"))
	}

	for _, allergen := range req.Allergens {
		if allergen != "" {
			database.DB.Create(&models.UserAllergen{
				UserID:   userID,
				Allergen: allergen,
			})
		}
	}

	var allergens []models.UserAllergen
	database.DB.Where("user_id = ?", userID).Find(&allergens)

	return c.JSON(models.Success(allergens))
}
