package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"uniqueIndex;not null;size:50"`
	Password  string         `json:"-" gorm:"not null;size:255"`
	Name      string         `json:"name" gorm:"not null;size:50"`
	Role      string         `json:"role" gorm:"type:enum('admin','chef','employee');not null;default:'employee'"`
	Avatar    string         `json:"avatar,omitempty" gorm:"size:500"`
	Allergens []UserAllergen `json:"allergens,omitempty" gorm:"foreignKey:UserID"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type UserAllergen struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	UserID    uint      `json:"user_id" gorm:"uniqueIndex:idx_user_allergen"`
	Allergen  string    `json:"allergen" gorm:"uniqueIndex:idx_user_allergen;not null;size:50"`
	CreatedAt time.Time `json:"created_at"`
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UpdateAllergensRequest struct {
	Allergens []string `json:"allergens"`
}
