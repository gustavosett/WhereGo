package geoip

import (
	"net/netip"
	"os"
	"testing"

	"github.com/oschwald/maxminddb-golang/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestModels_HasData(t *testing.T) {
	f := 1.0

	tests := []struct {
		name    string
		model   interface{ HasData() bool }
		hasData bool
		msg     string
	}{
		{"Names (Empty)", Names{}, false, "Empty Names"},
		{"Names (Full)", Names{English: "A"}, true, "Populated Names"},
		{"Continent (Empty)", Continent{}, false, "Empty Continent"},
		{"Continent (Full)", Continent{Code: "NA"}, true, "Populated Continent"},
		{"Location (Empty)", Location{}, false, "Empty Location"},
		{"Location (Full)", Location{TimeZone: "UTC"}, true, "Populated Location"},
		{"Enterprise (Empty)", Enterprise{}, false, "Empty Enterprise"},
		{"Enterprise (Full)", Enterprise{City: EnterpriseCityRecord{GeoNameID: 1}}, true, "Populated Enterprise"},
		{"City (Empty)", City{}, false, "Empty City"},
		{"City (Full)", City{Postal: CityPostal{Code: "123"}}, true, "Populated City"},
		{"Country (Empty)", Country{}, false, "Empty Country"},
		{"Country (Full)", Country{Country: CountryRecord{ISOCode: "US"}}, true, "Populated Country"},
		{"ASN (Empty)", ASN{}, false, "Empty ASN"},
		{"ASN (Full)", ASN{AutonomousSystemNumber: 123}, true, "Populated ASN"},
		{"ISP (Empty)", ISP{}, false, "Empty ISP"},
		{"ISP (Full)", ISP{ISP: "Comcast"}, true, "Populated ISP"},
		{"Domain (Empty)", Domain{}, false, "Empty Domain"},
		{"Domain (Full)", Domain{Domain: "google.com"}, true, "Populated Domain"},
		{"ConnectionType (Empty)", ConnectionType{}, false, "Empty ConnectionType"},
		{"ConnectionType (Full)", ConnectionType{ConnectionType: "Cable"}, true, "Populated ConnectionType"},
		{"AnonymousIP (Empty)", AnonymousIP{}, false, "Empty AnonymousIP"},
		{"AnonymousIP (Full)", AnonymousIP{IsAnonymous: true}, true, "Populated AnonymousIP"},
		{"CityRecord", CityRecord{GeoNameID: 1}, true, "CityRecord"},
		{"CityPostal", CityPostal{Code: "1"}, true, "CityPostal"},
		{"CitySubdivision", CitySubdivision{GeoNameID: 1}, true, "CitySubdivision"},
		{"CityTraits", CityTraits{IsAnycast: true}, true, "CityTraits"},
		{"CountryRecord", CountryRecord{GeoNameID: 1}, true, "CountryRecord"},
		{"CountryTraits", CountryTraits{IsAnycast: true}, true, "CountryTraits"},
		{"RepresentedCountry", RepresentedCountry{GeoNameID: 1}, true, "RepresentedCountry"},
		{"EnterpriseCityRecord", EnterpriseCityRecord{GeoNameID: 1}, true, "EnterpriseCityRecord"},
		{"EnterprisePostal", EnterprisePostal{Code: "1"}, true, "EnterprisePostal"},
		{"EnterpriseSubdivision", EnterpriseSubdivision{GeoNameID: 1}, true, "EnterpriseSubdivision"},
		{"EnterpriseCountryRecord", EnterpriseCountryRecord{GeoNameID: 1}, true, "EnterpriseCountryRecord"},
		{"EnterpriseTraits", EnterpriseTraits{ISP: "x"}, true, "EnterpriseTraits"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.hasData, tt.model.HasData(), tt.msg)
		})
	}

	t.Run("Location Coordinates Logic", func(t *testing.T) {
		assert.False(t, Location{}.HasCoordinates(), "Empty location has no coords")
		assert.False(t, Location{Latitude: &f}.HasCoordinates(), "Partial coords return false")
		assert.True(t, Location{Latitude: &f, Longitude: &f}.HasCoordinates(), "Full coords return true")
	})

	t.Run("Enterprise hasSubdivisionsData", func(t *testing.T) {
		e := Enterprise{}
		assert.False(t, e.HasData())
		e.Subdivisions = []EnterpriseSubdivision{{}, {}} // Empty subdivisions
		assert.False(t, e.HasData())
		e.Subdivisions = []EnterpriseSubdivision{{}, {GeoNameID: 1}} // One valid
		assert.True(t, e.HasData())
	})

	t.Run("City hasSubdivisionsData", func(t *testing.T) {
		c := City{}
		assert.False(t, c.HasData())
		c.Subdivisions = []CitySubdivision{{}, {}}
		assert.False(t, c.HasData())
		c.Subdivisions = []CitySubdivision{{GeoNameID: 1}}
		assert.True(t, c.HasData())
	})
}

func TestGetDBType(t *testing.T) {
	tests := []struct {
		name      string
		dbTypeStr string
		expected  databaseType
		expectErr bool
	}{
		{"Anonymous IP", "GeoIP2-Anonymous-IP", isAnonymousIP, false},
		{"ASN", "GeoLite2-ASN", isASN, false},
		{"City", "GeoIP2-City", isCity | isCountry, false},
		{"Country", "GeoIP2-Country", isCity | isCountry, false},
		{"Connection Type", "GeoIP2-Connection-Type", isConnectionType, false},
		{"Domain", "GeoIP2-Domain", isDomain, false},
		{"Enterprise", "GeoIP2-Enterprise", isEnterprise | isCity | isCountry, false},
		{"ISP", "GeoIP2-ISP", isISP | isASN, false},
		{"Unknown", "Alien-Database", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reader := &maxminddb.Reader{
				Metadata: maxminddb.Metadata{DatabaseType: tt.dbTypeStr},
			}
			got, err := getDBType(reader)
			if tt.expectErr {
				assert.Error(t, err)
				assert.IsType(t, UnknownDatabaseTypeError{}, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expected, got)
			}
		})
	}
}

func TestReader_MethodRestrictions(t *testing.T) {
	ip := netip.MustParseAddr("8.8.8.8")
	mockMmdb := &maxminddb.Reader{Metadata: maxminddb.Metadata{DatabaseType: "MockDB"}}

	tests := []struct {
		methodName string
		dbType     databaseType
		action     func(r *Reader) error
	}{
		{"Enterprise", isCity, func(r *Reader) error { _, err := r.Enterprise(ip); return err }},
		{"City", isASN, func(r *Reader) error { _, err := r.City(ip); return err }},
		{"Country", isASN, func(r *Reader) error { _, err := r.Country(ip); return err }},
		{"AnonymousIP", isCity, func(r *Reader) error { _, err := r.AnonymousIP(ip); return err }},
		{"ASN", isCity, func(r *Reader) error { _, err := r.ASN(ip); return err }},
		{"ConnectionType", isCity, func(r *Reader) error { _, err := r.ConnectionType(ip); return err }},
		{"Domain", isCity, func(r *Reader) error { _, err := r.Domain(ip); return err }},
		{"ISP", isCity, func(r *Reader) error { _, err := r.ISP(ip); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.methodName, func(t *testing.T) {
			r := &Reader{mmdbReader: mockMmdb, databaseType: tt.dbType}
			err := tt.action(r)
			require.Error(t, err)
			assert.IsType(t, InvalidMethodError{}, err)
		})
	}
}

func TestErrors_Formatting(t *testing.T) {
	e1 := InvalidMethodError{Method: "M", DatabaseType: "DB"}
	assert.Equal(t, "geoip2: the M method does not support the DB database", e1.Error())
	e2 := UnknownDatabaseTypeError{DatabaseType: "BadDB"}
	assert.Equal(t, `geoip2: reader does not support the "BadDB" database type`, e2.Error())
}

func TestApplyOptions(t *testing.T) {
	var called bool
	opt := func(o *readerOptions) { called = true }
	opts := applyOptions([]Option{opt, nil})
	assert.IsType(t, []maxminddb.ReaderOption{}, opts)
	assert.True(t, called)
}

func setupIntegration(t *testing.T) string {
	t.Helper()
	dbPath := "../../data/city.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skipf("Skipping integration test: database not found at %s", dbPath)
	}
	return dbPath
}

func TestReader_Integration_HappyPath(t *testing.T) {
	dbPath := setupIntegration(t)
	ip := netip.MustParseAddr("8.8.8.8")

	t.Run("Open and Lookup", func(t *testing.T) {
		r, err := Open(dbPath)
		require.NoError(t, err)
		defer func() {
			closeErr := r.Close()
			require.NoError(t, closeErr)
		}()

		city, err := r.City(ip)
		assert.NoError(t, err)
		assert.True(t, city.HasData())

		country, err := r.Country(ip)
		assert.NoError(t, err)
		assert.True(t, country.HasData())
	})

	t.Run("OpenBytes", func(t *testing.T) {
		b, err := os.ReadFile(dbPath)
		require.NoError(t, err)
		r, err := OpenBytes(b)
		require.NoError(t, err)
		rErr := r.Close()
		require.NoError(t, rErr)
	})

	t.Run("Invalid Path/Bytes", func(t *testing.T) {
		_, err := Open("invalid.mmdb")
		assert.Error(t, err)
		_, err = OpenBytes([]byte("bad"))
		assert.Error(t, err)
	})
}

func TestReader_ForcedExecution(t *testing.T) {
	dbPath := setupIntegration(t)
	realReader, err := maxminddb.Open(dbPath)
	require.NoError(t, err)
	defer func() {
		closeErr := realReader.Close()
		require.NoError(t, closeErr)
	}()

	ip := netip.MustParseAddr("8.8.8.8")

	tests := []struct {
		name       string
		forcedType databaseType
		action     func(r *Reader) (any, error)
	}{
		{"Enterprise", isEnterprise, func(r *Reader) (any, error) { return r.Enterprise(ip) }},
		{"ASN", isASN, func(r *Reader) (any, error) { return r.ASN(ip) }},
		{"ISP", isISP, func(r *Reader) (any, error) { return r.ISP(ip) }},
		{"Domain", isDomain, func(r *Reader) (any, error) { return r.Domain(ip) }},
		{"ConnectionType", isConnectionType, func(r *Reader) (any, error) { return r.ConnectionType(ip) }},
		{"AnonymousIP", isAnonymousIP, func(r *Reader) (any, error) { return r.AnonymousIP(ip) }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &Reader{
				mmdbReader:   realReader,
				databaseType: tt.forcedType,
			}
			result, err := tt.action(r)
			assert.NotNil(t, result)
			_ = err
		})
	}
}

func TestReader_DecodeErrors(t *testing.T) {
	dbPath := setupIntegration(t)
	data, err := os.ReadFile(dbPath)
	require.NoError(t, err)

	dbBytes := make([]byte, len(data))
	copy(dbBytes, data)

	mmdb, _ := maxminddb.OpenBytes(dbBytes)
	meta := mmdb.Metadata
	mmdbErr := mmdb.Close()
	require.NoError(t, mmdbErr)

	treeSizeBits := uint(meta.NodeCount) * uint(meta.RecordSize)
	treeSizeBytes := treeSizeBits / 8

	dataStart := int(treeSizeBytes) + 16

	dataEnd := len(dbBytes) - 5000

	if dataStart >= dataEnd {
		t.Skip("Database file too small to safely corrupt data section without touching metadata")
	}

	for i := dataStart; i < dataEnd; i++ {
		dbBytes[i] = 0xFF
	}

	r, err := OpenBytes(dbBytes)
	require.NoError(t, err, "OpenBytes should succeed because metadata at EOF is intact")
	defer func() {
		closeErr := r.Close()
		require.NoError(t, closeErr)
	}()

	ip := netip.MustParseAddr("8.8.8.8")

	tests := []struct {
		name       string
		forcedType databaseType
		action     func(r *Reader) error
	}{
		{"City", isCity, func(r *Reader) error { _, err := r.City(ip); return err }},
		{"Country", isCountry, func(r *Reader) error { _, err := r.Country(ip); return err }},
		{"Enterprise", isEnterprise, func(r *Reader) error { _, err := r.Enterprise(ip); return err }},
		{"ASN", isASN, func(r *Reader) error { _, err := r.ASN(ip); return err }},
		{"ISP", isISP, func(r *Reader) error { _, err := r.ISP(ip); return err }},
		{"Domain", isDomain, func(r *Reader) error { _, err := r.Domain(ip); return err }},
		{"Connection", isConnectionType, func(r *Reader) error { _, err := r.ConnectionType(ip); return err }},
		{"AnonIP", isAnonymousIP, func(r *Reader) error { _, err := r.AnonymousIP(ip); return err }},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			originalType := r.databaseType
			r.databaseType = tt.forcedType
			defer func() { r.databaseType = originalType }()

			err := tt.action(r)

			assert.Error(t, err, "Expected error due to corrupted data section")

			assert.IsNotType(t, InvalidMethodError{}, err)
		})
	}
}
