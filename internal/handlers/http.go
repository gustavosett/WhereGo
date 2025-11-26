package handlers

import (
	"net/http"

	"github.com/gustavosett/WhereGo/internal/geoip"
	"github.com/labstack/echo/v4"
)

type GeoIPHandler struct {
	GeoService *geoip.Service
}

var (
	errInvalidIP = map[string]string{"error": "invalid IP address"}
	errNoData    = map[string]string{"error": "no data found for the given IP"}
	healthOK     = map[string]string{"status": "ok"}
)

func (h *GeoIPHandler) Lookup(c echo.Context) error {
	result, err := h.GeoService.LookupIP(c.Param("ip"))
	if err != nil {
		if err == geoip.ErrInvalidIP {
			return c.JSON(http.StatusBadRequest, errInvalidIP)
		}
		return c.JSON(http.StatusNotFound, errNoData)
	}
	return c.JSON(http.StatusOK, result)
}

func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, healthOK)
}
