package handlers

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestHealthCheck(t *testing.T) {
	e := echo.New()
	e.GET("/health", HealthCheck)

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", rec.Code)
	}

	body, err := io.ReadAll(rec.Body)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	if !strings.Contains(string(body), `"status":"ok"`) {
		t.Errorf("Expected body to contain status ok, got '%s'", string(body))
	}

	contentType := rec.Header().Get("Content-Type")
	if !strings.HasPrefix(contentType, "application/json") {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestLookupIntegration(t *testing.T) {
	dbPath := "../../data/city.db"
	service, err := geoip.NewService(dbPath)
	if err != nil {
		t.Skipf("Skipping integration test: could not open database at %s: %v", dbPath, err)
	}
	defer func() {
		closeErr := service.DB.Close()
		require.NoError(t, closeErr)
	}()

	h := &GeoIPHandler{GeoService: service}
	e := echo.New()
	e.GET("/lookup/:ip", h.Lookup)

	// Test valid IP
	req := httptest.NewRequest(http.MethodGet, "/lookup/8.8.8.8", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("Expected status 200 for valid IP, got %d", rec.Code)
	}

	var result geoip.City
	if err := json.Unmarshal(rec.Body.Bytes(), &result); err != nil {
		t.Errorf("Failed to unmarshal response: %v", err)
	}

	// Test invalid IP
	req = httptest.NewRequest(http.MethodGet, "/lookup/invalid-ip", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid IP, got %d", rec.Code)
	}
}

func TestLookupDBError(t *testing.T) {
	dbPath := "../../data/city.db"
	service, err := geoip.NewService(dbPath)
	if err != nil {
		t.Skipf("Skipping integration test: could not open database at %s: %v", dbPath, err)
	}
	// Close immediately to simulate error
	closeErr := service.DB.Close()
	require.NoError(t, closeErr)

	h := &GeoIPHandler{GeoService: service}
	e := echo.New()
	e.GET("/lookup/:ip", h.Lookup)

	req := httptest.NewRequest(http.MethodGet, "/lookup/8.8.8.8", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for DB error, got %d", rec.Code)
	}
}
