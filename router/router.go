package router

import (
	"crypto/subtle"
	"log"
	"strings"

	"dday-backend/controllers/api"
	"dday-backend/controllers/rest"
	"dday-backend/global/config"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func SetupRoutes(app *fiber.App) {
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.New(cors.Config{
		AllowOrigins: corsOrigins(),
		AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders: "Origin,Content-Type,Accept,Authorization,X-API-Key",
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

	if config.AppConfig.Server.APIKey == "" {
		log.Println("WARNING: API_KEY is not set - write endpoints (POST/PUT/DELETE) are open to anyone. Set API_KEY in production.")
	}

	apiV1 := app.Group("/api/v1", apiKeyRequired)
	setupAPIRoutes(apiV1)

	rest := app.Group("/rest", apiKeyRequired)
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

// corsOrigins 는 CORS_ALLOWED_ORIGINS(쉼표 구분)를 돌려준다. 미설정이면 "*" 로 두고 경고한다.
func corsOrigins() string {
	var origins []string
	for _, o := range strings.Split(config.AppConfig.Server.CORSAllowedOrigins, ",") {
		if o = strings.TrimSpace(o); o != "" {
			origins = append(origins, o)
		}
	}
	if len(origins) == 0 {
		log.Println("WARNING: CORS_ALLOWED_ORIGINS is not set - allowing all origins (*). Set it in production.")
		return "*"
	}
	return strings.Join(origins, ",")
}

// apiKeyRequired 는 쓰기 요청(POST/PUT/DELETE/PATCH)에 X-API-Key 헤더를 요구한다.
// API_KEY 미설정 시에는 검사하지 않는다(개발용 — 기동 시 경고 로그).
// 주의: 공개 웹/앱 클라이언트에 넣은 키는 번들에서 추출 가능하므로 약한 장벽일 뿐 사용자 인증을 대신하지 못한다.
func apiKeyRequired(c *fiber.Ctx) error {
	switch c.Method() {
	case fiber.MethodGet, fiber.MethodHead, fiber.MethodOptions:
		return c.Next()
	}

	expected := config.AppConfig.Server.APIKey
	if expected == "" {
		return c.Next()
	}

	if subtle.ConstantTimeCompare([]byte(c.Get("X-API-Key")), []byte(expected)) != 1 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "invalid or missing API key",
		})
	}
	return c.Next()
}
