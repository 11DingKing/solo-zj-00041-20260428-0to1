package main

import (
	"fmt"
	"log"

	"canteen-system/config"
	"canteen-system/database"
	"canteen-system/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	config.InitConfig()

	database.InitDB()
	database.InitRedis()

	app := fiber.New(fiber.Config{
		AppName: "Canteen System API",
	})

	app.Use(recover.New())
	app.Use(logger.New(logger.Config{
		Format: "[${time}] ${status} - ${method} ${path}\n",
	}))

	routes.SetupRoutes(app)

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"message": "Canteen System API is running",
		})
	})

	port := config.AppConfig.Server.Port
	fmt.Printf("Server starting on port %s...\n", port)

	if err := app.Listen(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
