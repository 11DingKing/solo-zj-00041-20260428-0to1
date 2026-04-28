package handlers

import (
	"context"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"canteen-system/database"
	"canteen-system/models"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func generateOrderNo() string {
	now := time.Now()
	rand.Seed(now.UnixNano())
	datePart := now.Format("20060102150405")
	randomPart := fmt.Sprintf("%06d", rand.Intn(1000000))
	return "ORD" + datePart + randomPart
}

func CreateOrder(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)

	var req models.CreateOrderRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	if len(req.Items) == 0 {
		return c.Status(400).JSON(models.Error(400, "订单不能为空"))
	}

	var dailyMenu models.DailyMenu
	if err := database.DB.Where("menu_date = ? AND meal_period = ?", req.MenuDate, req.MealPeriod).
		Preload("Dishes.Dish").
		First(&dailyMenu).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "该时段菜单不存在"))
	}

	menuDishMap := make(map[uint]*models.DailyMenuDish)
	for i := range dailyMenu.Dishes {
		menuDishMap[dailyMenu.Dishes[i].ID] = &dailyMenu.Dishes[i]
	}

	var totalAmount float64
	var orderItems []models.OrderItem

	ctx := context.Background()
	pipe := database.Redis.Pipeline()

	for _, item := range req.Items {
		menuDish, exists := menuDishMap[item.DailyMenuDishID]
		if !exists {
			return c.Status(400).JSON(models.Error(400, "菜品不在菜单中"))
		}

		if menuDish.Dish == nil {
			return c.Status(400).JSON(models.Error(400, "菜品信息错误"))
		}

		stockKey := fmt.Sprintf("stock:%s:%s:%d", req.MenuDate, req.MealPeriod, menuDish.Dish.ID)
		pipe.DecrBy(ctx, stockKey, int64(item.Quantity))

		subtotal := menuDish.Dish.Price * float64(item.Quantity)
		totalAmount += subtotal

		orderItems = append(orderItems, models.OrderItem{
			DishID:    menuDish.Dish.ID,
			DishName:  menuDish.Dish.Name,
			DishPrice: menuDish.Dish.Price,
			Quantity:  item.Quantity,
			Subtotal:  subtotal,
		})
	}

	results, err := pipe.Exec(ctx)
	if err != nil {
		return c.Status(500).JSON(models.Error(500, "库存检查失败"))
	}

	rollbackNeeded := false
	for i, result := range results {
		cmd, ok := result.(*redis.IntCmd)
		if !ok {
			continue
		}
		newStock, _ := cmd.Result()
		if newStock < 0 {
			rollbackNeeded = true

			rollbackPipe := database.Redis.Pipeline()
			for j, item := range req.Items {
				if j <= i {
					menuDish := menuDishMap[item.DailyMenuDishID]
					stockKey := fmt.Sprintf("stock:%s:%s:%d", req.MenuDate, req.MealPeriod, menuDish.Dish.ID)
					rollbackPipe.IncrBy(ctx, stockKey, int64(item.Quantity))
				}
			}
			rollbackPipe.Exec(ctx)
			break
		}
	}

	if rollbackNeeded {
		return c.Status(400).JSON(models.Error(400, "部分菜品库存不足"))
	}

	pickupStart, err := time.Parse("2006-01-02 15:04:05", req.PickupTimeStart)
	if err != nil {
		pickupStart = time.Now()
	}
	pickupEnd, err := time.Parse("2006-01-02 15:04:05", req.PickupTimeEnd)
	if err != nil {
		pickupEnd = pickupStart.Add(time.Hour)
	}

	orderNo := generateOrderNo()
	order := models.Order{
		OrderNo:         orderNo,
		UserID:          userID,
		TotalAmount:     totalAmount,
		Status:          models.OrderStatusPendingConfirm,
		PickupTimeStart: pickupStart,
		PickupTimeEnd:   pickupEnd,
		QRCodeContent:   orderNo,
	}

	tx := database.DB.Begin()
	if err := tx.Create(&order).Error; err != nil {
		tx.Rollback()

		rollbackPipe := database.Redis.Pipeline()
		for _, item := range req.Items {
			menuDish := menuDishMap[item.DailyMenuDishID]
			stockKey := fmt.Sprintf("stock:%s:%s:%d", req.MenuDate, req.MealPeriod, menuDish.Dish.ID)
			rollbackPipe.IncrBy(ctx, stockKey, int64(item.Quantity))
		}
		rollbackPipe.Exec(ctx)

		return c.Status(500).JSON(models.Error(500, "创建订单失败"))
	}

	for i := range orderItems {
		orderItems[i].OrderID = order.ID
	}
	if err := tx.Create(&orderItems).Error; err != nil {
		tx.Rollback()

		rollbackPipe := database.Redis.Pipeline()
		for _, item := range req.Items {
			menuDish := menuDishMap[item.DailyMenuDishID]
			stockKey := fmt.Sprintf("stock:%s:%s:%d", req.MenuDate, req.MealPeriod, menuDish.Dish.ID)
			rollbackPipe.IncrBy(ctx, stockKey, int64(item.Quantity))
		}
		rollbackPipe.Exec(ctx)

		return c.Status(500).JSON(models.Error(500, "创建订单失败"))
	}

	for _, item := range req.Items {
		menuDish := menuDishMap[item.DailyMenuDishID]
		tx.Model(&models.DailyMenuDish{}).
			Where("id = ?", menuDish.ID).
			Update("remaining_quantity", gorm.Expr("remaining_quantity - ?", item.Quantity))
	}

	tx.Commit()

	database.DB.Preload("Items").First(&order, order.ID)

	return c.JSON(models.Success(order))
}

func GetMyOrders(c *fiber.Ctx) error {
	userID := c.Locals("user_id").(uint)
	status := c.Query("status")

	var orders []models.Order
	query := database.DB.Model(&models.Order{}).
		Preload("Items").
		Where("user_id = ?", userID)

	if status != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Order("created_at desc").Find(&orders).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取订单失败"))
	}

	return c.JSON(models.Success(orders))
}

func GetOrderByID(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的订单ID"))
	}

	userID := c.Locals("user_id").(uint)
	userRole := c.Locals("role").(string)

	var order models.Order
	query := database.DB.Model(&models.Order{}).
		Preload("Items").
		Preload("User")

	if userRole == "employee" {
		query = query.Where("id = ? AND user_id = ?", id, userID)
	} else {
		query = query.Where("id = ?", id)
	}

	if err := query.First(&order).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "订单不存在"))
	}

	return c.JSON(models.Success(order))
}

func GetOrderByNo(c *fiber.Ctx) error {
	orderNo := c.Params("order_no")
	if orderNo == "" {
		return c.Status(400).JSON(models.Error(400, "订单号不能为空"))
	}

	var order models.Order
	if err := database.DB.Model(&models.Order{}).
		Preload("Items").
		Preload("User").
		Where("order_no = ?", orderNo).
		First(&order).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "订单不存在"))
	}

	return c.JSON(models.Success(order))
}

func GetChefOrders(c *fiber.Ctx) error {
	status := c.Query("status")
	menuDate := c.Query("menu_date")
	mealPeriod := c.Query("meal_period")

	if menuDate == "" {
		menuDate = time.Now().Format("2006-01-02")
	}

	var orders []models.Order
	query := database.DB.Model(&models.Order{}).
		Preload("Items").
		Preload("User")

	if status != "" {
		query = query.Where("status = ?", status)
	} else {
		query = query.Where("status IN ?", []string{
			string(models.OrderStatusPendingConfirm),
			string(models.OrderStatusInProduction),
			string(models.OrderStatusReadyForPickup),
		})
	}

	query = query.Where("DATE(pickup_time_start) = ?", menuDate)

	if mealPeriod != "" {
		query = query.Where("meal_period = ?", mealPeriod)
	}

	if err := query.Order("pickup_time_start, created_at").Find(&orders).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "获取订单失败"))
	}

	return c.JSON(models.Success(orders))
}

func UpdateOrderStatus(c *fiber.Ctx) error {
	id, err := strconv.ParseUint(c.Params("id"), 10, 64)
	if err != nil {
		return c.Status(400).JSON(models.Error(400, "无效的订单ID"))
	}

	var req models.UpdateOrderStatusRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(400).JSON(models.Error(400, "请求参数错误"))
	}

	var order models.Order
	if err := database.DB.First(&order, id).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "订单不存在"))
	}

	validTransitions := map[models.OrderStatus][]models.OrderStatus{
		models.OrderStatusPendingConfirm:   {models.OrderStatusInProduction},
		models.OrderStatusInProduction:     {models.OrderStatusReadyForPickup},
		models.OrderStatusReadyForPickup:   {models.OrderStatusPickedUp},
		models.OrderStatusPickedUp:         {models.OrderStatusReviewed},
	}

	validNext := validTransitions[order.Status]
	isValid := false
	for _, s := range validNext {
		if s == req.Status {
			isValid = true
			break
		}
	}

	if !isValid {
		return c.Status(400).JSON(models.Error(400, "无效的订单状态转换"))
	}

	if err := database.DB.Model(&order).Update("status", req.Status).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "更新订单状态失败"))
	}

	return c.JSON(models.Success(nil))
}

func ConfirmPickup(c *fiber.Ctx) error {
	orderNo := c.Query("order_no")
	if orderNo == "" {
		return c.Status(400).JSON(models.Error(400, "订单号不能为空"))
	}

	var order models.Order
	if err := database.DB.Where("order_no = ?", orderNo).First(&order).Error; err != nil {
		return c.Status(404).JSON(models.Error(404, "订单不存在"))
	}

	if order.Status != models.OrderStatusReadyForPickup {
		return c.Status(400).JSON(models.Error(400, "订单状态不正确，无法取餐"))
	}

	if err := database.DB.Model(&order).Update("status", models.OrderStatusPickedUp).Error; err != nil {
		return c.Status(500).JSON(models.Error(500, "确认取餐失败"))
	}

	return c.JSON(models.Success(nil))
}
