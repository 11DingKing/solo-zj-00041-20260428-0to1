package routes

import (
	"canteen-system/handlers"
	"canteen-system/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func SetupRoutes(app *fiber.App) {
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
		AllowMethods: "GET, POST, PUT, DELETE, OPTIONS",
	}))

	api := app.Group("/api")

	auth := api.Group("/auth")
	auth.Post("/login", handlers.Login)

	api.Use(middleware.JWTAuth())

	api.Get("/user/me", handlers.GetCurrentUser)
	api.Put("/user/allergens", handlers.UpdateUserAllergens)

	api.Get("/categories", handlers.GetCategories)
	api.Post("/categories", middleware.RequireRole("admin"), handlers.CreateCategory)

	api.Get("/dishes", handlers.GetDishes)
	api.Get("/dishes/:id", handlers.GetDishByID)
	api.Post("/dishes", middleware.RequireRole("admin"), handlers.CreateDish)
	api.Put("/dishes/:id", middleware.RequireRole("admin"), handlers.UpdateDish)
	api.Delete("/dishes/:id", middleware.RequireRole("admin"), handlers.DeleteDish)

	api.Get("/menus", handlers.GetDailyMenus)
	api.Get("/menus/detail", handlers.GetDailyMenuDetail)
	api.Get("/menus/dates", handlers.GetAvailableDates)
	api.Post("/menus", middleware.RequireRole("admin"), handlers.CreateDailyMenu)
	api.Put("/menus/:id/dishes", middleware.RequireRole("admin"), handlers.UpdateDailyMenuDishes)

	api.Post("/menus/check-allergens", handlers.CheckAllergens)

	api.Get("/templates", middleware.RequireRole("admin"), handlers.GetWeeklyTemplates)
	api.Post("/templates", middleware.RequireRole("admin"), handlers.CreateWeeklyTemplate)
	api.Put("/templates/:id/items", middleware.RequireRole("admin"), handlers.UpdateTemplateItems)
	api.Post("/templates/:id/apply", middleware.RequireRole("admin"), handlers.ApplyWeeklyTemplate)

	api.Get("/orders", handlers.GetMyOrders)
	api.Get("/orders/:id", handlers.GetOrderByID)
	api.Post("/orders", handlers.CreateOrder)

	chef := api.Group("/chef")
	chef.Use(middleware.RequireRole("chef", "admin"))
	chef.Get("/orders", handlers.GetChefOrders)
	chef.Get("/orders/no/:order_no", handlers.GetOrderByNo)
	chef.Put("/orders/:id/status", handlers.UpdateOrderStatus)
	chef.Post("/pickup", handlers.ConfirmPickup)

	api.Get("/user/stats", handlers.GetUserStats)

	api.Post("/reviews/order/:order_id", handlers.CreateReview)
	api.Get("/reviews/dish/:dish_id", handlers.GetDishReviews)

	admin := api.Group("/admin")
	admin.Use(middleware.RequireRole("admin"))
	admin.Get("/dashboard", handlers.GetDashboardStats)
}
