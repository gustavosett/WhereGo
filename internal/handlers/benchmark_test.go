package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/labstack/echo/v4"
)

const benchDBPath = "../../data/city.db"

func BenchmarkLookup(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip("Database not available")
	}
	defer func() {
		if err := service.DB.Close(); err != nil {
			b.Logf("Failed to close database: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/lookup/:ip", (&GeoIPHandler{GeoService: service}).Lookup)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/lookup/8.8.8.8", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}
}

func BenchmarkLookupParallel(b *testing.B) {
	service, err := geoip.NewService(benchDBPath)
	if err != nil {
		b.Skip("Database not available")
	}
	defer func() {
		if err := service.DB.Close(); err != nil {
			b.Logf("Failed to close database: %v", err)
		}
	}()

	e := echo.New()
	e.GET("/lookup/:ip", (&GeoIPHandler{GeoService: service}).Lookup)

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			req := httptest.NewRequest(http.MethodGet, "/lookup/8.8.8.8", nil)
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
		}
	})
}

func BenchmarkHealth(b *testing.B) {
	e := echo.New()
	e.GET("/health", HealthCheck)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
	}
}
