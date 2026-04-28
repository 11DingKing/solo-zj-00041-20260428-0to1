package models

import (
	"encoding/json"
	"time"
)

type JSONString []string

func (j JSONString) MarshalJSON() ([]byte, error) {
	return json.Marshal([]string(j))
}

func (j *JSONString) UnmarshalJSON(data []byte) error {
	var s []string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	*j = s
	return nil
}

type OrderStatus string

const (
	OrderStatusPendingConfirm OrderStatus = "pending_confirm"
	OrderStatusInProduction   OrderStatus = "in_production"
	OrderStatusReadyForPickup OrderStatus = "ready_for_pickup"
	OrderStatusPickedUp       OrderStatus = "picked_up"
	OrderStatusReviewed       OrderStatus = "reviewed"
)

type Order struct {
	ID               uint        `json:"id" gorm:"primaryKey"`
	OrderNo          string      `json:"order_no" gorm:"uniqueIndex;not null;size:32"`
	UserID           uint        `json:"user_id" gorm:"not null;index"`
	User             *User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
	TotalAmount      float64     `json:"total_amount" gorm:"type:decimal(10,2);not null;default:0"`
	Status           OrderStatus `json:"status" gorm:"type:enum('pending_confirm','in_production','ready_for_pickup','picked_up','reviewed');not null;default:'pending_confirm';index"`
	PickupTimeStart  time.Time   `json:"pickup_time_start" gorm:"not null;index"`
	PickupTimeEnd    time.Time   `json:"pickup_time_end" gorm:"not null"`
	PickupCode       string      `json:"pickup_code,omitempty" gorm:"size:6"`
	QRCodeContent    string      `json:"qr_code_content,omitempty" gorm:"type:text"`
	Items            []OrderItem `json:"items,omitempty" gorm:"foreignKey:OrderID"`
	CreatedAt        time.Time   `json:"created_at" gorm:"index"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type OrderItem struct {
	ID         uint    `json:"id" gorm:"primaryKey"`
	OrderID    uint    `json:"order_id" gorm:"not null;index"`
	DishID     uint    `json:"dish_id" gorm:"not null;index"`
	DishName   string  `json:"dish_name" gorm:"not null;size:100"`
	DishPrice  float64 `json:"dish_price" gorm:"type:decimal(10,2);not null"`
	Quantity   int     `json:"quantity" gorm:"not null"`
	Subtotal   float64 `json:"subtotal" gorm:"type:decimal(10,2);not null"`
	CreatedAt  time.Time `json:"created_at"`
}

type Review struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	OrderID   uint      `json:"order_id" gorm:"not null;uniqueIndex:idx_order_dish_user"`
	UserID    uint      `json:"user_id" gorm:"not null;uniqueIndex:idx_order_dish_user"`
	DishID    uint      `json:"dish_id" gorm:"not null;uniqueIndex:idx_order_dish_user;index"`
	Dish      *Dish     `json:"dish,omitempty" gorm:"foreignKey:DishID"`
	Rating    int       `json:"rating" gorm:"not null;index"`
	Comment   string    `json:"comment,omitempty" gorm:"type:text"`
	CreatedAt time.Time `json:"created_at"`
}

type CreateOrderRequest struct {
	MenuDate         string             `json:"menu_date" validate:"required"`
	MealPeriod       string             `json:"meal_period" validate:"required"`
	Items            []OrderItemRequest `json:"items" validate:"required,min=1"`
	PickupTimeStart  string             `json:"pickup_time_start" validate:"required"`
	PickupTimeEnd    string             `json:"pickup_time_end" validate:"required"`
}

type OrderItemRequest struct {
	DailyMenuDishID uint `json:"daily_menu_dish_id" validate:"required"`
	Quantity         int  `json:"quantity" validate:"required,min=1"`
}

type UpdateOrderStatusRequest struct {
	Status OrderStatus `json:"status" validate:"required"`
}

type CreateReviewRequest struct {
	Rating  int    `json:"rating" validate:"required,min=1,max=5"`
	Comment string `json:"comment"`
}

type AllergenCheckRequest struct {
	DailyMenuDishIDs []uint `json:"daily_menu_dish_ids" validate:"required"`
}

type AllergenCheckResponse struct {
	HasAllergens   bool              `json:"has_allergens"`
	AllergenDishes []AllergenDishInfo `json:"allergen_dishes"`
}

type AllergenDishInfo struct {
	DishID     uint     `json:"dish_id"`
	DishName   string   `json:"dish_name"`
	Allergens  []string `json:"allergens"`
}

type OrderListRequest struct {
	Status     string `form:"status"`
	MenuDate   string `form:"menu_date"`
	MealPeriod string `form:"meal_period"`
	Page       int    `form:"page,default=1"`
	PageSize   int    `form:"page_size,default=20"`
}

type OrderListResponse struct {
	Total    int64   `json:"total"`
	Page     int     `json:"page"`
	PageSize int     `json:"page_size"`
	Orders   []Order `json:"orders"`
}
