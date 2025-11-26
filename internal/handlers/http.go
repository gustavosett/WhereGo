package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
)

type GeoIPHandler struct {
	GeoService *geoip.Service
}

type Response struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	City     string `json:"city"`
	ISOCode  string `json:"iso_code"`
	Timezone string `json:"timezone"`
}

func (h *GeoIPHandler) Lookup(c *fiber.Ctx) error {
	ipParam := c.Params("ip")

	data, err := h.GeoService.LookupIP(ipParam)

	if err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	if data == nil {
		return c.Status(fiber.StatusNotFound).SendString("No data found for the given IP")
	}

	response := Response{
		IP:       ipParam,
		Country:  data.Country,
		City:     data.City,
		ISOCode:  data.ISOCode,
		Timezone: data.Timezone,
	}

	return c.JSON(response)
}
