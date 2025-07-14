package router

import (
	"dday-backend/controllers/api"
	"dday-backend/controllers/rest"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization",
	}))

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "D-Day Backend API",
			"version": "2.0",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	apiV1 := app.Group("/api/v1")
	setupAPIRoutes(apiV1)

	rest := app.Group("/rest")
	setupRESTRoutes(rest)
}

func setupAPIRoutes(router fiber.Router) {
	ddayAPI := &api.DdayController{}

	ddays := router.Group("/ddays")
	ddays.Get("/", ddayAPI.GetDdays)
	ddays.Post("/", ddayAPI.CreateDday)
	ddays.Get("/:id", ddayAPI.GetDday)
	ddays.Put("/:id", ddayAPI.UpdateDday)
	ddays.Delete("/:id", ddayAPI.DeleteDday)
}

func setupRESTRoutes(router fiber.Router) {
	ddayREST := &rest.DdayController{}

	ddays := router.Group("/ddays")
	ddays.Get("/", ddayREST.List)
	ddays.Post("/", ddayREST.Create)
	ddays.Get("/:id", ddayREST.Get)
	ddays.Put("/:id", ddayREST.Update)
	ddays.Delete("/:id", ddayREST.Delete)
}