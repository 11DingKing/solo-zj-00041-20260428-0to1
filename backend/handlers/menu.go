package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"canteen-system/database"
	"canteen-system/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

func getMenuCacheKey(date, mealPeriod string) string {
	return fmt.Sprintf("menu:%s:%s", date, mealPeriod)
}

func GetDailyMenus(c *fiber.Ctx) error {
	date := c.Query("date")
	mealPeriod := c.Query("meal_period")

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	var menus []models.DailyMenu
	query := database.DB.Model(&models.DailyMenu{}).
		Preload("Dishes.Dish.Category").
		Where("menu_date = ?", date).
		Where("is_active = ?", true)

	if mealPeriod != "" {
		query = query.Where("meal_period = ?", mealPeriod)
	}

	if err := query.Order("FIELD(meal_period, 'breakfast', 'lunch', 'dinner')").Find(&menus).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取菜单失败"))
	}

	return c.JSON(models.Success(menus))
}

func GetDailyMenuDetail(c *fiber.Ctx) error {
	date := c.Query("date")
	mealPeriod := c.Query("meal_period")

	if date == "" || mealPeriod == "" {
		return c.Status(400).JSON(models.Error(400, "日期和时段不能为空"))
	}

	cacheKey := getMenuCacheKey(date, mealPeriod)
	ctx := context.Background()

	cached, err := database.Redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var menu models.DailyMenu
		if json.Unmarshal([]byte(cached), &menu) == nil {
			return c.JSON(models.Success(menu))
		}
	}

	var menu models.DailyMenu
	if err := database.DB.Model(&models.DailyMenu{}).
		Preload("Dishes.Dish.Category").
		Where("menu_date = ? AND meal_period = ?", date, mealPeriod).
		First(&menu).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "菜单不存在"))
	}

	menuJSON, _ := json.Marshal(menu)
	database.Redis.Set(ctx, cacheKey, menuJSON, 30*time.Minute)

	return c.JSON(models.Success(menu))
}

func CreateDailyMenu(c *fiber.Ctx) error {
	var req models.CreateDailyMenuRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if req.MenuDate == "" || req.MealPeriod == "" {
		return c.Status(400).JSON(models.Error(400, "日期和时段不能为空"))
	}

	var existingMenu models.DailyMenu
	if err := database.DB.Where("menu_date = ? AND meal_period = ?", req.MenuDate, req.MealPeriod).First(&existingMenu).Error; err == nil {
		return c.Status(400).JSON(models.Error(400, "该时段菜单已存在"))
	}

	var dishes []models.Dish
	if err := database.DB.Where("id IN ?", req.DishIDs).Find(&dishes).Error; err != nil {
		return c.Status(400).JSON(models.Error(400, "部分菜品不存在"))
	}

	if len(dishes) != len(req.DishIDs) {
		return c.Status(400).JSON(models.Error(400, "部分菜品不存在"))
	}

	menu := models.DailyMenu{
		MenuDate:   req.MenuDate,
		MealPeriod: req.MealPeriod,
		StartTime:  req.StartTime,
		EndTime:    req.EndTime,
		IsActive:   true,
	}

	if err := database.DB.Create(&menu).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "创建菜单失败"))
	}

	for _, dish := range dishes {
		menuDish := models.DailyMenuDish{
			DailyMenuID:       menu.ID,
			DishID:            dish.ID,
			RemainingQuantity: dish.DailyLimit,
		}
		database.DB.Create(&menuDish)

		stockKey := fmt.Sprintf("stock:%s:%s:%d", req.MenuDate, req.MealPeriod, dish.ID)
		database.Redis.Set(context.Background(), stockKey, dish.DailyLimit, 24*time.Hour)
	}

	ctx := context.Background()
	cacheKey := getMenuCacheKey(req.MenuDate, req.MealPeriod)
	database.Redis.Del(ctx, cacheKey)

	return c.JSON(models.Success(menu))
}

func UpdateDailyMenuDishes(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的菜单ID"))
	}

	var menu models.DailyMenu
	if err := database.DB.First(&menu, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "菜单不存在"))
	}

	var req struct {
		DishIDs []uint `json:"dish_ids" validate:"required"`
	}
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if err := database.DB.Where("daily_menu_id = ?", menu.ID).Delete(&models.DailyMenuDish{}).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "更新菜品失败"))
	}

	var dishes []models.Dish
	database.DB.Where("id IN ?", req.DishIDs).Find(&dishes)

	for _, dish := range dishes {
		menuDish := models.DailyMenuDish{
			DailyMenuID:       menu.ID,
			DishID:            dish.ID,
			RemainingQuantity: dish.DailyLimit,
		}
		database.DB.Create(&menuDish)

		stockKey := fmt.Sprintf("stock:%s:%s:%d", menu.MenuDate, menu.MealPeriod, dish.ID)
		database.Redis.Set(context.Background(), stockKey, dish.DailyLimit, 24*time.Hour)
	}

	ctx := context.Background()
	cacheKey := getMenuCacheKey(menu.MenuDate, menu.MealPeriod)
	database.Redis.Del(ctx, cacheKey)

	return c.JSON(models.Success(nil))
}

func GetWeeklyTemplates(c *fiber.Ctx) error {
	var templates []models.WeeklyMenuTemplate
	if err := database.DB.Preload("Items.Dish").Order("id desc").Find(&templates).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取模板失败"))
	}
	return c.JSON(models.Success(templates))
}

func CreateWeeklyTemplate(c *fiber.Ctx) error {
	var req models.CreateTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	template := models.WeeklyMenuTemplate{
		Name:        req.Name,
		Description: req.Description,
		IsActive:    false,
	}

	if err := database.DB.Create(&template).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "创建模板失败"))
	}

	return c.JSON(models.Success(template))
}

func UpdateTemplateItems(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的模板ID"))
	}

	var template models.WeeklyMenuTemplate
	if err := database.DB.First(&template, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "模板不存在"))
	}

	var items []models.TemplateItem
	if err := c.BodyParser(&items); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if err := database.DB.Where("template_id = ?", template.ID).Delete(&models.WeeklyMenuTemplateItem{}).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "更新模板失败"))
	}

	for _, item := range items {
		for _, dishID := range item.DishIDs {
			templateItem := models.WeeklyMenuTemplateItem{
				TemplateID: template.ID,
				DayOfWeek:  item.DayOfWeek,
				MealPeriod: item.MealPeriod,
				DishID:     dishID,
			}
			database.DB.Create(&templateItem)
		}
	}

	return c.JSON(models.Success(nil))
}

func ApplyWeeklyTemplate(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的模板ID"))
	}

	var req models.ApplyTemplateRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	var template models.WeeklyMenuTemplate
	if err := database.DB.Preload("Items.Dish").First(&template, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "模板不存在"))
	}

	startDate, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的日期格式"))
	}

	startDayOfWeek := int(startDate.Weekday())
	if startDayOfWeek == 0 {
		startDayOfWeek = 7
	}

	itemsByDay := make(map[int]map[string][]models.Dish)
	for _, item := range template.Items {
		if itemsByDay[item.DayOfWeek] == nil {
			itemsByDay[item.DayOfWeek] = make(map[string][]models.Dish)
		}
		if item.Dish != nil {
			itemsByDay[item.DayOfWeek][item.MealPeriod] = append(itemsByDay[item.DayOfWeek][item.MealPeriod], *item.Dish)
		}
	}

	for dayOffset := 0; dayOffset < 7; dayOffset++ {
		currentDate := startDate.AddDate(0, 0, dayOffset)
		currentDayOfWeek := (startDayOfWeek + dayOffset - 1) % 7 + 1

		dayItems, exists := itemsByDay[currentDayOfWeek]
		if !exists {
			continue
		}

		defaultTimes := map[string][2]string{
			"breakfast": {"07:00:00", "09:00:00"},
			"lunch":     {"11:00:00", "13:00:00"},
			"dinner":    {"17:00:00", "19:00:00"},
		}

		for mealPeriod, dishes := range dayItems {
			times := defaultTimes[mealPeriod]

			var existingMenu models.DailyMenu
			if err := database.DB.Where("menu_date = ? AND meal_period = ?", currentDate.Format("2006-01-02"), mealPeriod).First(&existingMenu).Error; err == nil {
				continue
			}

			menu := models.DailyMenu{
				MenuDate:   currentDate.Format("2006-01-02"),
				MealPeriod: mealPeriod,
				StartTime:  times[0],
				EndTime:    times[1],
				IsActive:   true,
			}

			if err := database.DB.Create(&menu).Error; err != nil {
				continue
			}

			for _, dish := range dishes {
				menuDish := models.DailyMenuDish{
					DailyMenuID:       menu.ID,
					DishID:            dish.ID,
					RemainingQuantity: dish.DailyLimit,
				}
				database.DB.Create(&menuDish)

				stockKey := fmt.Sprintf("stock:%s:%s:%d", menu.MenuDate, mealPeriod, dish.ID)
				database.Redis.Set(context.Background(), stockKey, dish.DailyLimit, 24*time.Hour)
			}
		}
	}

	return c.JSON(models.Success(nil))
}

func GetAvailableDates(c *fiber.Ctx) error {
	today := time.Now()
	tomorrow := today.AddDate(0, 0, 1)

	dates := []map[string]interface{}{
		{
			"date":    today.Format("2006-01-02"),
			"label":   "今天",
			"weekday": today.Weekday().String(),
		},
		{
			"date":    tomorrow.Format("2006-01-02"),
			"label":   "明天",
			"weekday": tomorrow.Weekday().String(),
		},
	}

	return c.JSON(models.Success(dates))
}

func CheckAllergens(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req models.AllergenCheckRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	var userAllergens []models.UserAllergen
	database.DB.Where("user_id = ?", userID).Find(&userAllergens)

	if len(userAllergens) == 0 {
		return c.JSON(models.Success(models.AllergenCheckResponse{
			HasAllergens:   false,
			AllergenDishes: []models.AllergenDishInfo{},
		}))
	}

	userAllergenSet := make(map[string]bool)
	for _, a := range userAllergens {
		userAllergenSet[a.Allergen] = true
	}

	var menuDishes []models.DailyMenuDish
	database.DB.Preload("Dish").Where("id IN ?", req.DailyMenuDishIDs).Find(&menuDishes)

	var allergenDishes []models.AllergenDishInfo
	for _, md := range menuDishes {
		if md.Dish == nil {
			continue
		}

		dishAllergens := md.Dish.Allergens
		hasMatch := false
		matchingAllergens := []string{}

		for _, a := range dishAllergens {
			if userAllergenSet[a] {
				hasMatch = true
				matchingAllergens = append(matchingAllergens, a)
			}
		}

		if hasMatch {
			allergenDishes = append(allergenDishes, models.AllergenDishInfo{
				DishID:    md.Dish.ID,
				DishName:  md.Dish.Name,
				Allergens: matchingAllergens,
			})
		}
	}

	return c.JSON(models.Success(models.AllergenCheckResponse{
		HasAllergens:   len(allergenDishes) > 0,
		AllergenDishes: allergenDishes,
	}))
}
