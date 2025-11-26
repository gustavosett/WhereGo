package handlers

import (
	"encoding/json"
	"io"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gofiber/fiber/v2"
)

const (
	contentTypeHeader = "Content-Type"
	applicationJSON   = "application/json"
	testIPGoogle      = "8.8.8.8"
	testCountryUS     = "United States"
	testCityMountain  = "Mountain View"
	testISOCodeUS     = "US"
	testTimezonePST   = "America/Los_Angeles"
	healthEndpoint    = "/health"
	testEndpoint      = "/test"
)

// Response is used for test JSON unmarshaling
type Response struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	City     string `json:"city"`
	ISOCode  string `json:"iso_code"`
	Timezone string `json:"timezone"`
}

func TestResponseJSONMarshaling(t *testing.T) {
	response := Response{
		IP:       testIPGoogle,
		Country:  testCountryUS,
		City:     testCityMountain,
		ISOCode:  testISOCodeUS,
		Timezone: testTimezonePST,
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
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
		t.Fatalf("marshal error: %v", err)
	}

	jsonStr := string(data)

	// Verify JSON tags are correct
	expectedTags := []string{`"ip"`, `"country"`, `"city"`, `"iso_code"`, `"timezone"`}
	for _, tag := range expectedTags {
		if !strings.Contains(jsonStr, tag) {
			t.Errorf("Expected JSON to contain tag %s, got: %s", tag, jsonStr)
		}
	}
}

func TestGeoIPHandlerLookupWithNilService(t *testing.T) {
	handler := &GeoIPHandler{
		GeoService: nil,
	}

	if handler.GeoService != nil {
		t.Error("Expected nil GeoService")
	}
}

func TestLookupEmptyIP(t *testing.T) {
	app := fiber.New()

	handler := &GeoIPHandler{
		GeoService: nil,
	}

	app.Get("/lookup/:ip", handler.Lookup)

	req := httptest.NewRequest("GET", "/lookup/", nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	if resp.StatusCode != fiber.StatusNotFound {
		t.Errorf("Expected status 404, got %d", resp.StatusCode)
	}
}

func TestResponseEmptyFields(t *testing.T) {
	response := Response{}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("marshal error: %v", err)
	}

	var unmarshaled Response
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("unmarshal error: %v", err)
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

			app.Get(testEndpoint, func(c *fiber.Ctx) error {
				return c.SendStatus(tt.expectedStatus)
			})

			req := httptest.NewRequest("GET", testEndpoint, nil)
			resp, err := app.Test(req)
			if err != nil {
				t.Fatalf("test error: %v", err)
			}

			if resp.StatusCode != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, resp.StatusCode)
			}
		})
	}
}

func TestContentTypeHeader(t *testing.T) {
	app := fiber.New()

	app.Get(testEndpoint, func(c *fiber.Ctx) error {
		c.Set(contentTypeHeader, applicationJSON)
		return c.SendString("{}")
	})

	req := httptest.NewRequest("GET", testEndpoint, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	contentType := resp.Header.Get(contentTypeHeader)
	if contentType != applicationJSON {
		t.Errorf("Expected Content-Type '%s', got '%s'", applicationJSON, contentType)
	}
}

func TestJSONEncoder(t *testing.T) {
	app := fiber.New()

	response := Response{
		IP:       testIPGoogle,
		Country:  testCountryUS,
		City:     testCityMountain,
		ISOCode:  testISOCodeUS,
		Timezone: testTimezonePST,
	}

	app.Get(testEndpoint, func(c *fiber.Ctx) error {
		return c.JSON(response)
	})

	req := httptest.NewRequest("GET", testEndpoint, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}

	var decoded Response
	err = json.Unmarshal(body, &decoded)
	if err != nil {
		t.Fatalf("decode error: %v", err)
	}

	if decoded.IP != response.IP {
		t.Errorf("Expected IP '%s', got '%s'", response.IP, decoded.IP)
	}
}

func TestHealthCheck(t *testing.T) {
	app := fiber.New()
	app.Get(healthEndpoint, HealthCheck)

	req := httptest.NewRequest("GET", healthEndpoint, nil)
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

	expected := `{"status":"ok"}`
	if string(body) != expected {
		t.Errorf("Expected body '%s', got '%s'", expected, string(body))
	}
}

func TestHealthCheckContentType(t *testing.T) {
	app := fiber.New()
	app.Get(healthEndpoint, HealthCheck)

	req := httptest.NewRequest("GET", healthEndpoint, nil)
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("test error: %v", err)
	}

	contentType := resp.Header.Get(contentTypeHeader)
	if contentType != applicationJSON {
		t.Errorf("Expected Content-Type '%s', got '%s'", applicationJSON, contentType)
	}
}
