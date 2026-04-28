package handlers

import (
	"strconv"

	"canteen-system/database"
	"canteen-system/models"

	"github.com/gofiber/fiber/v2"
)

func CreateReview(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	orderID, err := strconv.ParseUint(c.Params("order_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的订单ID"))
	}

	var order models.Order
	if err := database.DB.Preload("Items").First(&order, orderID).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "订单不存在"))
	}

	if order.UserID != userID {
		return c.Status(403).JSON(models.Error(403, "无权评价此订单"))
	}

	if order.Status != models.OrderStatusPickedUp && order.Status != models.OrderStatusReviewed {
		return c.Status(400).JSON(models.Error(400, "订单未完成，无法评价"))
	}

	var reviews []map[string]interface{}
	if err := c.BodyParser(&reviews); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	for _, reviewData := range reviews {
		dishID, ok := reviewData["dish_id"].(float64)
		if !ok {
			continue
		}
		rating, ok := reviewData["rating"].(float64)
		if !ok || rating < 1 || rating > 5 {
			continue
		}
		comment, _ := reviewData["comment"].(string)

		found := false
		for _, item := range order.Items {
			if item.DishID == uint(dishID) {
				found = true
				break
			}
		}
		if !found {
			continue
		}

		var existingReview models.Review
		if err := database.DB.Where("order_id = ? AND dish_id = ? AND user_id = ?", orderID, uint(dishID), userID).First(&existingReview).Error; err == nil {
			existingReview.Rating = int(rating)
			existingReview.Comment = comment
			database.DB.Save(&existingReview)
			continue
		}

		review := models.Review{
			OrderID: uint(orderID),
			UserID:  userID,
			DishID:  uint(dishID),
			Rating:  int(rating),
			Comment: comment,
		}
		database.DB.Create(&review)
	}

	hasUnreviewed := false
	for _, item := range order.Items {
		var count int64
		database.DB.Model(&models.Review{}).Where("order_id = ? AND dish_id = ?", orderID, item.DishID).Count(&count)
		if count == 0 {
			hasUnreviewed = true
			break
		}
	}

	if !hasUnreviewed && order.Status == models.OrderStatusPickedUp {
		database.DB.Model(&order).Update("status", models.OrderStatusReviewed)
	}

	return c.JSON(models.Success(nil))
}

func GetDishReviews(c *fiber.Ctx) error {
	dishID, err := strconv.ParseUint(c.Params("dish_id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的菜品ID"))
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	pageSize, _ := strconv.Atoi(c.Query("page_size", "10"))

	offset := (page - 1) * pageSize

	var total int64
	database.DB.Model(&models.Review{}).Where("dish_id = ?", dishID).Count(&total)

	var reviews []models.Review
	if err := database.DB.Preload("Dish").
		Where("dish_id = ?", dishID).
		Order("created_at desc").
		Offset(offset).
		Limit(pageSize).
		Find(&reviews).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取评价失败"))
	}

	return c.JSON(models.Success(map[string]interface{}{
		"total":     total,
		"page":      page,
		"page_size": pageSize,
		"reviews":   reviews,
	}))
}

func GetUserStats(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	currentTime := "2026-04-28"

	var monthAmount float64
	database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("user_id = ? AND DATE_FORMAT(created_at, '%Y-%m') = ?", userID, currentTime[:7]).
		Scan(&monthAmount)

	var orderCount int64
	database.DB.Model(&models.Order{}).
		Where("user_id = ? AND DATE_FORMAT(created_at, '%Y-%m') = ?", userID, currentTime[:7]).
		Count(&orderCount)

	type TopDishResult struct {
		DishID   uint
		DishName string
		Count    int64
	}

	var topDishes []TopDishResult
	database.DB.Raw(`
		SELECT oi.dish_id, oi.dish_name, SUM(oi.quantity) as count
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE o.user_id = ? AND DATE_FORMAT(o.created_at, '%Y-%m') = ?
		GROUP BY oi.dish_id, oi.dish_name
		ORDER BY count DESC
		LIMIT 5
	`, userID, currentTime[:7]).Scan(&topDishes)

	var userTopDishes []models.UserTopDish
	for _, td := range topDishes {
		userTopDishes = append(userTopDishes, models.UserTopDish{
			DishID:   td.DishID,
			DishName: td.DishName,
			Count:    td.Count,
		})
	}

	return c.JSON(models.Success(models.UserStats{
		MonthAmount: monthAmount,
		OrderCount:  orderCount,
		TopDishes:   userTopDishes,
	}))
}

func GetDashboardStats(c *fiber.Ctx) error {
	today := "2026-04-28"
	tomorrow := "2026-04-29"

	var todayOrders int64
	database.DB.Model(&models.Order{}).
		Where("DATE(created_at) = ?", today).
		Count(&todayOrders)

	var todayRevenue float64
	database.DB.Model(&models.Order{}).
		Select("COALESCE(SUM(total_amount), 0)").
		Where("DATE(created_at) = ?", today).
		Scan(&todayRevenue)

	type MealPeriodResult struct {
		MealPeriod string
		Count      int64
	}

	var mealPeriodDistribution []MealPeriodResult
	database.DB.Raw(`
		SELECT 
			CASE 
				WHEN HOUR(pickup_time_start) < 10 THEN 'breakfast'
				WHEN HOUR(pickup_time_start) < 16 THEN 'lunch'
				ELSE 'dinner'
			END as meal_period,
			COUNT(*) as count
		FROM orders
		WHERE DATE(created_at) = ?
		GROUP BY meal_period
	`, today).Scan(&mealPeriodDistribution)

	var periodCounts []models.MealPeriodCount
	for _, mp := range mealPeriodDistribution {
		periodCounts = append(periodCounts, models.MealPeriodCount{
			MealPeriod: mp.MealPeriod,
			Count:      mp.Count,
		})
	}

	type TopDishResult struct {
		DishID       uint
		DishName     string
		TotalOrders  int64
		TotalRevenue float64
	}

	var topDishes []TopDishResult
	database.DB.Raw(`
		SELECT 
			oi.dish_id, 
			oi.dish_name, 
			SUM(oi.quantity) as total_orders,
			SUM(oi.subtotal) as total_revenue
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE DATE(o.created_at) >= DATE_SUB(?, INTERVAL 30 DAY)
		GROUP BY oi.dish_id, oi.dish_name
		ORDER BY total_orders DESC
		LIMIT 10
	`, today).Scan(&topDishes)

	var dishTopList []models.TopDish
	for _, td := range topDishes {
		dishTopList = append(dishTopList, models.TopDish{
			DishID:       td.DishID,
			DishName:     td.DishName,
			TotalOrders:  td.TotalOrders,
			TotalRevenue: td.TotalRevenue,
		})
	}

	type DailyRevenueResult struct {
		Date    string
		Revenue float64
	}

	var last30DaysRevenue []DailyRevenueResult
	database.DB.Raw(`
		SELECT 
			DATE(created_at) as date,
			COALESCE(SUM(total_amount), 0) as revenue
		FROM orders
		WHERE DATE(created_at) >= DATE_SUB(?, INTERVAL 30 DAY)
		GROUP BY DATE(created_at)
		ORDER BY date
	`, today).Scan(&last30DaysRevenue)

	var dailyRevenueList []models.DailyRevenue
	for _, dr := range last30DaysRevenue {
		dailyRevenueList = append(dailyRevenueList, models.DailyRevenue{
			Date:    dr.Date,
			Revenue: dr.Revenue,
		})
	}

	type RatedDishResult struct {
		DishID      uint
		DishName    string
		AvgRating   float64
		ReviewCount int64
	}

	var topRatedDishes []RatedDishResult
	database.DB.Raw(`
		SELECT 
			r.dish_id,
			d.name as dish_name,
			AVG(r.rating) as avg_rating,
			COUNT(*) as review_count
		FROM reviews r
		JOIN dishes d ON r.dish_id = d.id
		GROUP BY r.dish_id, d.name
		HAVING review_count >= 5
		ORDER BY avg_rating DESC
		LIMIT 10
	`).Scan(&topRatedDishes)

	var ratedList []models.RatedDish
	for _, rd := range topRatedDishes {
		ratedList = append(ratedList, models.RatedDish{
			DishID:      rd.DishID,
			DishName:    rd.DishName,
			AvgRating:   rd.AvgRating,
			ReviewCount: rd.ReviewCount,
		})
	}

	type IngredientResult struct {
		DishID   uint
		DishName string
		Quantity int64
	}

	var tomorrowIngredients []IngredientResult
	database.DB.Raw(`
		SELECT 
			oi.dish_id,
			oi.dish_name,
			SUM(oi.quantity) as quantity
		FROM order_items oi
		JOIN orders o ON oi.order_id = o.id
		WHERE DATE(o.pickup_time_start) = ?
		GROUP BY oi.dish_id, oi.dish_name
		ORDER BY quantity DESC
	`, tomorrow).Scan(&tomorrowIngredients)

	var ingredientList []models.IngredientDemand
	for _, ig := range tomorrowIngredients {
		ingredientList = append(ingredientList, models.IngredientDemand{
			DishID:   ig.DishID,
			DishName: ig.DishName,
			Quantity: ig.Quantity,
		})
	}

	return c.JSON(models.Success(models.DashboardStats{
		TodayOrders:          todayOrders,
		TodayRevenue:         todayRevenue,
		MealPeriodDistribution: periodCounts,
		TopDishes:            dishTopList,
		Last30DaysRevenue:    dailyRevenueList,
		TopRatedDishes:       ratedList,
		TomorrowIngredients:  ingredientList,
	}))
}
