package geoip

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOpenInvalidPath(t *testing.T) {
	_, err := Open("nonexistent/path/to/database.mmdb")
	if err == nil {
		t.Error("Expected error when opening non-existent database, got nil")
	}
}

func TestNamesHasData(t *testing.T) {
	var emptyNames Names
	assert.False(t, emptyNames.HasData(), "Empty Names should not have data")

	nonEmptyNames := Names{English: "Test"}
	assert.True(t, nonEmptyNames.HasData(), "Names with data should have data")
}

func TestAllStructsHaveHasData(t *testing.T) {
	// Ensure all result structs have HasData methods
	var city City
	var country Country
	var enterprise Enterprise
	var anonymousIP AnonymousIP
	var asn ASN
	var connectionType ConnectionType
	var domain Domain
	var isp ISP
	var names Names

	// These should all compile and return false for zero values (no data)
	assert.False(t, city.HasData())
	assert.False(t, country.HasData())
	assert.False(t, enterprise.HasData())
	assert.False(t, anonymousIP.HasData())
	assert.False(t, asn.HasData())
	assert.False(t, connectionType.HasData())
	assert.False(t, domain.HasData())
	assert.False(t, isp.HasData())
	assert.False(t, names.HasData())
}

func TestASNHasData(t *testing.T) {
	var emptyASN ASN
	assert.False(t, emptyASN.HasData(), "Empty ASN should not have data")

	nonEmptyASN := ASN{AutonomousSystemNumber: 123}
	assert.True(t, nonEmptyASN.HasData(), "ASN with data should have data")
}

func TestLocationHasCoordinates(t *testing.T) {
	var emptyLocation Location
	assert.False(t, emptyLocation.HasCoordinates(), "Empty Location should not have coordinates")

	lat := 51.5142
	lon := -0.0931
	locationWithCoords := Location{Latitude: &lat, Longitude: &lon}
	assert.True(t, locationWithCoords.HasCoordinates(), "Location with lat/lon should have coordinates")

	// Only latitude
	locationOnlyLat := Location{Latitude: &lat}
	assert.False(t, locationOnlyLat.HasCoordinates(), "Location with only latitude should not have coordinates")
}
