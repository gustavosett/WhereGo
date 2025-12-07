package main

import (
	"errors"
	"log"
	"net/http"
	"os"

	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/gustavosett/WhereGo/internal/handlers"
	jsoniter "github.com/json-iterator/go"
	"github.com/labstack/echo/v4"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func NewServer(dbPath string) (*echo.Echo, *geoip.Service, error) {
	geoService, err := geoip.NewService(dbPath)
	if err != nil {
		return nil, nil, err
	}

	handler := &handlers.GeoIPHandler{
		GeoService: geoService,
	}

	e := echo.New()
	e.JSONSerializer = &JSONSerializer{}

	e.GET("/health", handlers.HealthCheck)
	e.GET("/lookup/:ip", handler.Lookup)

	return e, geoService, nil
}

func main() {
	e, geoService, err := NewServer("data/city.db")
	if err != nil {
		log.Fatalf("Failed to initialize GeoIP service: %v", err)
	}
	defer func() {
		if err := geoService.DB.Close(); err != nil {
			log.Printf("Failed to close GeoIP database: %v", err)
		}
	}()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on :%s", port)
	if err := e.Start(":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatalf("Server failed: %v", err)
	}
}

// JSONSerializer implements echo.JSONSerializer using json-iterator
type JSONSerializer struct{}

func (s *JSONSerializer) Serialize(c echo.Context, i interface{}, indent string) error {
	enc := json.NewEncoder(c.Response())
	if indent != "" {
		enc.SetIndent("", indent)
	}
	return enc.Encode(i)
}

func (s *JSONSerializer) Deserialize(c echo.Context, i interface{}) error {
	return json.NewDecoder(c.Request().Body).Decode(i)
}
