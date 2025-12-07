package geoip

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewService(t *testing.T) {
	t.Run("Invalid Path", func(t *testing.T) {
		svc, err := NewService("nonexistent/path/to/database.mmdb")
		assert.Error(t, err)
		assert.Nil(t, svc)
	})

	t.Run("Success Integration", func(t *testing.T) {
		dbPath := "../../data/city.db"
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			t.Skipf("Skipping integration test: database not found at %s", dbPath)
		}

		svc, err := NewService(dbPath)
		require.NoError(t, err)
		assert.NotNil(t, svc)
		assert.NotNil(t, svc.DB)

		// Cleanup
		svcErr := svc.DB.Close()
		require.NoError(t, svcErr)
	})
}

func TestLookupIP_Validation(t *testing.T) {
	svc := &Service{DB: nil}

	tests := []struct {
		name  string
		ipStr string
	}{
		{"Alphabetical", "invalid-ip"},
		{"Empty", ""},
		{"Partial IPv4", "192.168"},
		{"IPv4 with Port", "8.8.8.8:80"},
		{"IPv6 with Port", "[2001:db8::1]:80"},
		{"Out of Range", "300.300.300.300"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city, err := svc.LookupIP(tt.ipStr)

			assert.ErrorIs(t, err, ErrInvalidIP, "Should return sentinel error for invalid IP")
			assert.Nil(t, city)
		})
	}
}

func TestLookupIP_Integration(t *testing.T) {
	dbPath := "../../data/city.db"
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		t.Skipf("Skipping integration test: database not found at %s", dbPath)
	}

	svc, err := NewService(dbPath)
	require.NoError(t, err)
	defer func() {
		closeErr := svc.DB.Close()
		require.NoError(t, closeErr)
	}()

	tests := []struct {
		name        string
		ipStr       string
		expectFound bool
		expectedISO string
	}{
		{
			name:        "Valid Google DNS (US)",
			ipStr:       "8.8.8.8",
			expectFound: true,
			expectedISO: "US", // Usually resolves to US
		},
		{
			name:        "Valid Cloudflare DNS (IPv6)",
			ipStr:       "2606:4700:4700::1111",
			expectFound: true,
		},
		{
			name:        "Valid Localhost (No Data)",
			ipStr:       "127.0.0.1",
			expectFound: false, // Localhost isn't in GeoIP DBs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			city, err := svc.LookupIP(tt.ipStr)
			require.NoError(t, err, "Lookup should not return error for valid IP format")
			require.NotNil(t, city, "City struct should never be nil even if empty")

			if tt.expectFound {
				assert.True(t, city.HasData(), "Expected data for IP %s", tt.ipStr)
				if tt.expectedISO != "" {
					assert.Equal(t, tt.expectedISO, city.Country.ISOCode)
				}
			} else {
				assert.False(t, city.HasData(), "Expected no data for IP %s", tt.ipStr)
			}
		})
	}
}
