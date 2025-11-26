package handlers

import (
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get("/health", HealthCheck)

	req := httptest.NewRequest("GET", "/health", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	if resp.StatusCode != fiber.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if !strings.Contains(string(body), `"status":"ok"`) {
		t.Errorf("Expected body to contain status ok, got '%s'", string(body))
	}

	contentType := resp.Header.Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestLookupWithNilService(t *testing.T) {
	app := fiber.New()
	handler := &GeoIPHandler{GeoService: nil}
	app.Get("/lookup/:ip", handler.Lookup)

	// Route without IP parameter returns 404
	req := httptest.NewRequest("GET", "/lookup/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}
