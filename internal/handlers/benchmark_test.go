package handlers

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
)

const benchDBPath = "../../data/city.db"

func BenchmarkLookup(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip("Database not available")
	}
	defer service.DB.Close() //nolint:errcheck

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/lookup/:ip", (&GeoIPHandler{GeoService: service}).Lookup)

	req := httptest.NewRequest("GET", "/lookup/8.8.8.8", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close() //nolint:errcheck
	}
}

func BenchmarkLookupParallel(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip("Database not available")
	}
	defer service.DB.Close() //nolint:errcheck

	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/lookup/:ip", (&GeoIPHandler{GeoService: service}).Lookup)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", "/lookup/8.8.8.8", nil)
		for pb.Next() {
			resp, _ := app.Test(req, -1)
			resp.Body.Close() //nolint:errcheck
		}
	})
}

func BenchmarkHealth(b *testing.B) {
	app := fiber.New(fiber.Config{DisableStartupMessage: true})
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close() //nolint:errcheck
	}
}
