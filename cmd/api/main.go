package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/gustavosett/WhereGo/internal/handlers"
)

func main() {
	geoService, err := geoip.NewService("data/city.db")
	if err != nil {
		log.Fatalf("Failed to initialize GeoIP service: %v", err)
	}
	defer geoService.DB.Close()

	handler := &handlers.GeoIPHandler{
		GeoService: geoService,
	}

	app := fiber.New()
	app.Use(logger.New())
	app.Use(recover.New())

	app.Get("/health", handlers.HealthCheck)
	app.Get("/lookup/:ip", handler.Lookup)

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
