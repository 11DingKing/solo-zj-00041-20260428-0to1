package models

import (
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"uniqueIndex;not null;size:50"`
	SortOrder int            `json:"sort_order" gorm:"default:0"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

type Dish struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	Name        string         `json:"name" gorm:"not null;size:100"`
	CategoryID  uint           `json:"category_id" gorm:"not null;index"`
	Category    *Category      `json:"category,omitempty" gorm:"foreignKey:CategoryID"`
	Price       float64        `json:"price" gorm:"type:decimal(10,2);not null;default:0"`
	Image       string         `json:"image,omitempty" gorm:"size:500"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	DailyLimit  int            `json:"daily_limit" gorm:"not null;default:50"`
	Allergens   JSONString     `json:"allergens,omitempty" gorm:"type:json"`
	IsAvailable bool           `json:"is_available" gorm:"default:true"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `json:"-" gorm:"index"`
}

type WeeklyMenuTemplate struct {
	ID          uint                        `json:"id" gorm:"primaryKey"`
	Name        string                      `json:"name" gorm:"not null;size:100"`
	Description string                      `json:"description,omitempty" gorm:"type:text"`
	IsActive    bool                        `json:"is_active" gorm:"default:false"`
	Items       []WeeklyMenuTemplateItem    `json:"items,omitempty" gorm:"foreignKey:TemplateID"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type WeeklyMenuTemplateItem struct {
	ID          uint           `json:"id" gorm:"primaryKey"`
	TemplateID  uint           `json:"template_id" gorm:"not null;uniqueIndex:idx_template_day_period_dish;index"`
	DayOfWeek   int            `json:"day_of_week" gorm:"not null;uniqueIndex:idx_template_day_period_dish"`
	MealPeriod  string         `json:"meal_period" gorm:"type:enum('breakfast','lunch','dinner');not null;uniqueIndex:idx_template_day_period_dish"`
	DishID      uint           `json:"dish_id" gorm:"not null;uniqueIndex:idx_template_day_period_dish;index"`
	Dish        *Dish          `json:"dish,omitempty" gorm:"foreignKey:DishID"`
	CreatedAt   time.Time      `json:"created_at"`
}

type DailyMenu struct {
	ID         uint               `json:"id" gorm:"primaryKey"`
	MenuDate   string             `json:"menu_date" gorm:"type:date;not null;uniqueIndex:idx_date_period"`
	MealPeriod string             `json:"meal_period" gorm:"type:enum('breakfast','lunch','dinner');not null;uniqueIndex:idx_date_period"`
	StartTime  string             `json:"start_time" gorm:"type:time;not null"`
	EndTime    string             `json:"end_time" gorm:"type:time;not null"`
	IsActive   bool               `json:"is_active" gorm:"default:true"`
	Dishes     []DailyMenuDish    `json:"dishes,omitempty" gorm:"foreignKey:DailyMenuID"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
}

type DailyMenuDish struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	DailyMenuID      uint      `json:"daily_menu_id" gorm:"not null;uniqueIndex:idx_menu_dish;index"`
	DishID           uint      `json:"dish_id" gorm:"not null;uniqueIndex:idx_menu_dish;index"`
	Dish             *Dish     `json:"dish,omitempty" gorm:"foreignKey:DishID"`
	RemainingQuantity int      `json:"remaining_quantity" gorm:"not null"`
	CreatedAt        time.Time `json:"created_at"`
}

type CreateDishRequest struct {
	Name        string   `json:"name" validate:"required"`
	CategoryID  uint     `json:"category_id" validate:"required"`
	Price       float64  `json:"price" validate:"required,min=0"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	DailyLimit  int      `json:"daily_limit" validate:"min=1"`
	Allergens   []string `json:"allergens"`
	IsAvailable bool     `json:"is_available"`
}

type UpdateDishRequest struct {
	Name        string   `json:"name"`
	CategoryID  uint     `json:"category_id"`
	Price       float64  `json:"price" validate:"min=0"`
	Image       string   `json:"image"`
	Description string   `json:"description"`
	DailyLimit  int      `json:"daily_limit" validate:"min=1"`
	Allergens   []string `json:"allergens"`
	IsAvailable *bool    `json:"is_available"`
}

type CreateDailyMenuRequest struct {
	MenuDate   string `json:"menu_date" validate:"required"`
	MealPeriod string `json:"meal_period" validate:"required"`
	StartTime  string `json:"start_time" validate:"required"`
	EndTime    string `json:"end_time" validate:"required"`
	DishIDs    []uint `json:"dish_ids" validate:"required,min=1"`
}

type CreateTemplateRequest struct {
	Name        string `json:"name" validate:"required"`
	Description string `json:"description"`
}

type TemplateItem struct {
	DayOfWeek  int    `json:"day_of_week" validate:"required,min=1,max=7"`
	MealPeriod string `json:"meal_period" validate:"required"`
	DishIDs    []uint `json:"dish_ids"`
}

type ApplyTemplateRequest struct {
	StartDate string `json:"start_date" validate:"required"`
}
