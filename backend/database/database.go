package database

import (
	"fmt"
	"log"

	"canteen-system/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := config.AppConfig.Database.GetDSN()

	logLevel := logger.Info
	if config.AppConfig.Server.Mode == "release" {
		logLevel = logger.Silent
	}

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		log.Fatalf("Failed to get sql DB: %v", err)
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)

	fmt.Println("Database connected successfully!")
}
