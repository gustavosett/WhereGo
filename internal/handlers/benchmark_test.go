package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
)

const (
	benchDBPath       = "../../data/city.db"
	benchSkipMsg      = "Database not available for benchmarking"
	benchIPGoogle     = "8.8.8.8"
	benchIPCloudflare = "1.1.1.1"
	benchLookupRoute  = "/lookup/:ip"
	benchLookupPath   = "/lookup/"
)

//nolint:errcheck // benchmark cleanup
func BenchmarkLookupRoute(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip(benchSkipMsg)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	handler := &GeoIPHandler{GeoService: service}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get(benchLookupRoute, handler.Lookup)

	req := httptest.NewRequest("GET", benchLookupPath+benchIPGoogle, nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}

//nolint:errcheck // benchmark cleanup
func BenchmarkLookupRouteParallel(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip(benchSkipMsg)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	handler := &GeoIPHandler{GeoService: service}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get(benchLookupRoute, handler.Lookup)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		req := httptest.NewRequest("GET", benchLookupPath+benchIPGoogle, nil)
		for pb.Next() {
			resp, _ := app.Test(req, -1)
			resp.Body.Close()
		}
	})
}

//nolint:errcheck // benchmark cleanup
func BenchmarkLookupRouteMultipleIPs(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip(benchSkipMsg)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	handler := &GeoIPHandler{GeoService: service}

	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get(benchLookupRoute, handler.Lookup)

	ips := []string{
		benchIPGoogle,
		benchIPCloudflare,
		"208.67.222.222",
		"9.9.9.9",
		"185.228.168.9",
	}

	reqs := make([]*http.Request, len(ips))
	for i, ip := range ips {
		reqs[i] = httptest.NewRequest("GET", benchLookupPath+ip, nil)
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(reqs[i%len(reqs)], -1)
		resp.Body.Close()
	}
}

//nolint:errcheck // benchmark cleanup
func BenchmarkHealthRoute(b *testing.B) {
	app := fiber.New(fiber.Config{
		DisableStartupMessage: true,
	})
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		resp, _ := app.Test(req, -1)
		resp.Body.Close()
	}
}
