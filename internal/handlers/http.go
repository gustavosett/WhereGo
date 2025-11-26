package handlers

import (
	"unsafe"

	"github.com/gofiber/fiber/v2"
	"github.com/gustavosett/WhereGo/internal/geoip"
)

type GeoIPHandler struct {
	GeoService *geoip.Service
}

// Pre-allocated byte slices for zero-allocation responses
var (
	errInvalidIP     = []byte("invalid IP address")
	errNoData        = []byte("no data found for the given IP")
	healthOKResponse = []byte(`{"status":"ok"}`)
	jsonContentType  = []byte("application/json")

	// JSON fragments
	jsonStart    = []byte(`{"ip":"`)
	jsonCountry  = []byte(`","country":"`)
	jsonCity     = []byte(`","city":"`)
	jsonISOCode  = []byte(`","iso_code":"`)
	jsonTimezone = []byte(`","timezone":"`)
	jsonEnd      = []byte(`"}`)
)

// s2b converts string to []byte without allocation
func s2b(s string) []byte {
	return unsafe.Slice(unsafe.StringData(s), len(s))
}

//nolint:errcheck // buffer writes to response body don't fail
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

	c.Response().Header.SetContentTypeBytes(jsonContentType)
	buf := c.Response().BodyWriter()
	buf.Write(jsonStart)
	buf.Write(s2b(ipParam))
	buf.Write(jsonCountry)
	buf.Write(s2b(data.Country))
	buf.Write(jsonCity)
	buf.Write(s2b(data.City))
	buf.Write(jsonISOCode)
	buf.Write(s2b(data.ISOCode))
	buf.Write(jsonTimezone)
	buf.Write(s2b(data.Timezone))
	buf.Write(jsonEnd)
	return nil
}

func HealthCheck(c *fiber.Ctx) error {
	c.Response().Header.SetContentTypeBytes(jsonContentType)
	return c.Send(healthOKResponse)
}
