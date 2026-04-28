package models

type DashboardStats struct {
	TodayOrders          int64              `json:"today_orders"`
	TodayRevenue         float64            `json:"today_revenue"`
	MealPeriodDistribution []MealPeriodCount `json:"meal_period_distribution"`
	TopDishes            []TopDish          `json:"top_dishes"`
	Last30DaysRevenue    []DailyRevenue     `json:"last_30_days_revenue"`
	TopRatedDishes       []RatedDish        `json:"top_rated_dishes"`
	TomorrowIngredients  []IngredientDemand `json:"tomorrow_ingredients"`
}

type MealPeriodCount struct {
	MealPeriod string `json:"meal_period"`
	Count      int64  `json:"count"`
}

type TopDish struct {
	DishID      uint    `json:"dish_id"`
	DishName    string  `json:"dish_name"`
	TotalOrders int64   `json:"total_orders"`
	TotalRevenue float64 `json:"total_revenue"`
}

type DailyRevenue struct {
	Date    string  `json:"date"`
	Revenue float64 `json:"revenue"`
}

type RatedDish struct {
	DishID      uint    `json:"dish_id"`
	DishName    string  `json:"dish_name"`
	AvgRating   float64 `json:"avg_rating"`
	ReviewCount int64   `json:"review_count"`
}

type IngredientDemand struct {
	DishID       uint    `json:"dish_id"`
	DishName     string  `json:"dish_name"`
	Quantity     int64   `json:"quantity"`
}

type UserStats struct {
	MonthAmount   float64      `json:"month_amount"`
	OrderCount    int64        `json:"order_count"`
	TopDishes     []UserTopDish `json:"top_dishes"`
}

type UserTopDish struct {
	DishID     uint   `json:"dish_id"`
	DishName   string `json:"dish_name"`
	Count      int64  `json:"count"`
}

type ApiResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) ApiResponse {
	return ApiResponse{
		Code:    200,
		Message: "success",
		Data:    data,
	}
}

func Error(code int, message string) ApiResponse {
	return ApiResponse{
		Code:    code,
		Message: message,
	}
}
