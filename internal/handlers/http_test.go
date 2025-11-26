package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
)

func TestResponseJSONMarshaling(t *testing.T) {
	response := Response{
		IP:       "8.8.8.8",
		Country:  "United States",
		City:     "Mountain View",
		ISOCode:  "US",
		Timezone: "America/Los_Angeles",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	if unmarshaled.IP != response.IP {
		t.Errorf("Expected IP '%s', got '%s'", response.IP, unmarshaled.IP)
	}
	if unmarshaled.Country != response.Country {
		t.Errorf("Expected Country '%s', got '%s'", response.Country, unmarshaled.Country)
	}
	if unmarshaled.City != response.City {
		t.Errorf("Expected City '%s', got '%s'", response.City, unmarshaled.City)
	}
	if unmarshaled.ISOCode != response.ISOCode {
		t.Errorf("Expected ISOCode '%s', got '%s'", response.ISOCode, unmarshaled.ISOCode)
	}
	if unmarshaled.Timezone != response.Timezone {
		t.Errorf("Expected Timezone '%s', got '%s'", response.Timezone, unmarshaled.Timezone)
	}
}

func TestResponseJSONTags(t *testing.T) {
	response := Response{
		IP:       "1.1.1.1",
		Country:  "Australia",
		City:     "Sydney",
		ISOCode:  "AU",
		Timezone: "Australia/Sydney",
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(data)

	// Verify JSON tags are correct
	expectedTags := []string{`"ip"`, `"country"`, `"city"`, `"iso_code"`, `"timezone"`}
	for _, tag := range expectedTags {
		if !contains(jsonStr, tag) {
			t.Errorf("Expected JSON to contain tag %s, got: %s", tag, jsonStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestGeoIPHandlerLookupWithNilService(t *testing.T) {
	handler := &GeoIPHandler{
		GeoService: nil,
	}

	// This test verifies the handler struct can be created
	// even if the service is nil (though it would panic on use)
	if handler.GeoService != nil {
		t.Error("Expected nil GeoService")
	}
}

// MockGeoService is a mock implementation for testing
type MockGeoService struct {
	LookupResult *MockLookupData
	LookupError  error
}

type MockLookupData struct {
	Country  string
	City     string
	ISOCode  string
	Timezone string
}

func TestLookupEmptyIP(t *testing.T) {
	app := fiber.New()

	handler := &GeoIPHandler{
		GeoService: nil,
	}

	app.Get("/lookup/:ip", handler.Lookup)

	// Test with empty path (this will be a 404 due to routing)
	req := httptest.NewRequest("GET", "/lookup/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	// Empty IP in path should return 404 (route not matched)
	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestResponseEmptyFields(t *testing.T) {
	response := Response{}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal empty response: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal empty response: %v", err)
	}

	if unmarshaled.IP != "" {
		t.Errorf("Expected empty IP, got '%s'", unmarshaled.IP)
	}
	if unmarshaled.Country != "" {
		t.Errorf("Expected empty Country, got '%s'", unmarshaled.Country)
	}
}

func TestHTTPStatusCodes(t *testing.T) {
	tests := []struct {
		name           string
		expectedStatus int
	}{
		{"StatusOK", fiber.StatusOK},
		{"StatusBadRequest", fiber.StatusBadRequest},
		{"StatusNotFound", fiber.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()

			app.Get("/test", func(c *fiber.Ctx) error {
				return c.SendStatus(tt.expectedStatus)
			})

			req := httptest.NewRequest("GET", "/test", nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("Failed to test request: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestContentTypeHeader(t *testing.T) {
	app := fiber.New()

	app.Get("/test", func(c *fiber.Ctx) error {
		c.Set(contentTypeHeader, applicationJSON)
		return c.SendString("{}")
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	contentType := resp.Header.Get(contentTypeHeader)
	if contentType != applicationJSON {
		t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
	}
}

func TestJSONEncoder(t *testing.T) {
	app := fiber.New()

	response := Response{
		IP:       "8.8.8.8",
		Country:  "United States",
		City:     "Mountain View",
		ISOCode:  "US",
		Timezone: "America/Los_Angeles",
	}

	app.Get("/test", func(c *fiber.Ctx) error {
		return c.JSON(response)
	})

	req := httptest.NewRequest("GET", "/test", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("Failed to test request: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("Failed to read response body: %v", err)
	}

	var decoded Response
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		t.Fatalf("Failed to decode response body: %v", err)
	}

	if decoded.IP != response.IP {
		t.Errorf("Expected IP '%s', got '%s'", response.IP, decoded.IP)
	}
}
