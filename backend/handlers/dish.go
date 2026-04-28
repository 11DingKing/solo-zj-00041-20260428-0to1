package handlers

import (
	"strconv"

	"canteen-system/database"
	"canteen-system/models"

	"github.com/gofiber/fiber/v2"
)

func GetCategories(c *fiber.Ctx) error {
	var categories []models.Category
	if err := database.DB.Order("sort_order").Find(&categories).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取分类失败"))
	}
	return c.JSON(models.Success(categories))
}

func CreateCategory(c *fiber.Ctx) error {
	var category models.Category
	if err := c.BodyParser(&category); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if category.Name == "" {
		return c.Status(400).JSON(models.Error(400, "分类名称不能为空"))
	}

	if err := database.DB.Create(&category).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "创建分类失败"))
	}

	return c.JSON(models.Success(category))
}

func GetDishes(c *fiber.Ctx) error {
	categoryID := c.Query("category_id")
	isAvailable := c.Query("is_available")

	var dishes []models.Dish
	query := database.DB.Model(&models.Dish{}).Preload("Category")

	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if isAvailable != "" {
		avail, _ := strconv.ParseBool(isAvailable)
		query = query.Where("is_available = ?", avail)
	}

	if err := query.Order("id desc").Find(&dishes).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取菜品列表失败"))
	}

	return c.JSON(models.Success(dishes))
}

func GetDishByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的菜品ID"))
	}

	var dish models.Dish
	if err := database.DB.Preload("Category").First(&dish, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "菜品不存在"))
	}

	var reviews []models.Review
	database.DB.Where("dish_id = ?", dish.ID).Preload("Dish").Order("created_at desc").Limit(20).Find(&reviews)

	type DishWithReviews struct {
		models.Dish
		Reviews []models.Review `json:"reviews"`
	}

	return c.JSON(models.Success(DishWithReviews{
		Dish:    dish,
		Reviews: reviews,
	}))
}

func CreateDish(c *fiber.Ctx) error {
	var req models.CreateDishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if req.Name == "" || req.CategoryID == 0 {
		return c.Status(400).JSON(models.Error(400, "菜品名称和分类不能为空"))
	}

	var category models.Category
	if err := database.DB.First(&category, req.CategoryID).Error; err != nil {
		return c.Status(400).JSON(models.Error(400, "分类不存在"))
	}

	dish := models.Dish{
		Name:        req.Name,
		CategoryID:  req.CategoryID,
		Price:       req.Price,
		Image:       req.Image,
		Description: req.Description,
		DailyLimit:  req.DailyLimit,
		Allergens:   req.Allergens,
		IsAvailable: req.IsAvailable,
	}

	if dish.DailyLimit == 0 {
		dish.DailyLimit = 50
	}

	if err := database.DB.Create(&dish).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "创建菜品失败"))
	}

	return c.JSON(models.Success(dish))
}

func UpdateDish(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的菜品ID"))
	}

	var req models.UpdateDishRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	var dish models.Dish
	if err := database.DB.First(&dish, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "菜品不存在"))
	}

	if req.Name != "" {
		dish.Name = req.Name
	}
	if req.CategoryID > 0 {
		dish.CategoryID = req.CategoryID
	}
	if req.Price > 0 {
		dish.Price = req.Price
	}
	if req.Image != "" {
		dish.Image = req.Image
	}
	if req.Description != "" {
		dish.Description = req.Description
	}
	if req.DailyLimit > 0 {
		dish.DailyLimit = req.DailyLimit
	}
	if req.Allergens != nil {
		dish.Allergens = req.Allergens
	}
	if req.IsAvailable != nil {
		dish.IsAvailable = *req.IsAvailable
	}

	if err := database.DB.Save(&dish).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "更新菜品失败"))
	}

	return c.JSON(models.Success(dish))
}

func DeleteDish(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的菜品ID"))
	}

	if err := database.DB.Delete(&models.Dish{}, id).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "删除菜品失败"))
	}

	return c.JSON(models.Success(nil))
}
