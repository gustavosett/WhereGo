package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
)

type GeoIPHandler struct {
	GeoService *geoip.Service
}

var (
	errInvalidIP = fiber.Map{"error": "invalid IP address"}
	errNoData    = fiber.Map{"error": "no data found for the given IP"}
	healthOK     = fiber.Map{"status": "ok"}
)

func (h *GeoIPHandler) Lookup(c *fiber.Ctx) error {
	result, err := h.GeoService.LookupIP(c.Params("ip"))
	if err != nil {
		if err == geoip.ErrInvalidIP {
			return c.Status(fiber.StatusBadRequest).JSON(errInvalidIP)
		}
		return c.Status(fiber.StatusNotFound).JSON(errNoData)
	}
	return c.JSON(result)
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(healthOK)
}
