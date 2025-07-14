package main

import (
	"dday-backend/global/config"
	"dday-backend/models"
	"dday-backend/router"
	"log"

	"github.com/gofiber/fiber/v2"
)

func main() {
	config.LoadConfig()

	if err := models.InitDatabase(); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}
	defer models.DB.Close()

	app := fiber.New(fiber.Config{
		AppName: "D-Day Backend API v2.0",
	})

	router.SetupRoutes(app)

	port := ":" + config.AppConfig.Server.Port
	log.Printf("Server starting on port %s", port)
	log.Fatal(app.Listen(port))
}