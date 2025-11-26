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

var (
	errInvalidIP = []byte("invalid IP address")
	errNoData    = []byte("no data found for the given IP")
)

func (h *GeoIPHandler) Lookup(c *fiber.Ctx) error {
	ipParam := c.Params("ip")

	data := h.GeoService.GetLookupData()
	defer h.GeoService.PutLookupData(data)

	err := h.GeoService.LookupIP(ipParam, data)
	if err != nil {
		if err == geoip.ErrInvalidIP {
			return c.Status(fiber.StatusBadRequest).Send(errInvalidIP)
		}
		return c.Status(fiber.StatusNotFound).Send(errNoData)
	}

	return c.JSON(Response{
		IP:       ipParam,
		Country:  data.Country,
		City:     data.City,
		ISOCode:  data.ISOCode,
		Timezone: data.Timezone,
	})
}

type HealthResponse struct {
	Status string `json:"status"`
}

func HealthCheck(c *fiber.Ctx) error {
	return c.JSON(HealthResponse{Status: "ok"})
}
