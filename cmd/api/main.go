package main

import (
	"log"
	"os"
	"runtime"
	"time"

	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/gustavosett/WhereGo/internal/handlers"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	geoService, err := geoip.NewService("data/city.db")
	if err != nil {
		log.Fatalf("Failed to initialize GeoIP service: %v", err)
	}
	defer geoService.DB.Close() //nolint:errcheck

	handler := &handlers.GeoIPHandler{
		GeoService: geoService,
	}

	app := fiber.New(fiber.Config{
		JSONEncoder:               json.Marshal,
		JSONDecoder:               json.Unmarshal,
		DisableStartupMessage:     true,
		Prefork:                   os.Getenv("PREFORK") == "true",
		ReduceMemoryUsage:         true,
		StrictRouting:             true,
		CaseSensitive:             true,
		DisableDefaultDate:        true,
		DisableDefaultContentType: true,
		ReadTimeout:               5 * time.Second,
		WriteTimeout:              5 * time.Second,
		IdleTimeout:               120 * time.Second,
		ReadBufferSize:            4096,
		WriteBufferSize:           4096,
		CompressedFileSuffix:      ".gz",
		GETOnly:                   true,
	})

	app.Get("/health", handlers.HealthCheck)
	app.Get("/lookup/:ip", handler.Lookup)

	log.Println("Starting server on :8080")
	if err := app.Listen(":8080"); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
