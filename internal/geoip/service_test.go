package geoip

import (
	"net/netip"
	"testing"
)

func TestParseIP(t *testing.T) {
	tests := []struct {
		name    string
		ipStr   string
		isValid bool
	}{
		{"valid IPv4", "8.8.8.8", true},
		{"valid IPv6", "2001:4860:4860::8888", true},
		{"valid localhost", "127.0.0.1", true},
		{"invalid IP", "invalid", false},
		{"empty string", "", false},
		{"partial IP", "192.168", false},
		{"IP with port", "8.8.8.8:80", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := netip.ParseAddr(tt.ipStr)
			if tt.isValid && err != nil {
				t.Errorf("Expected valid IP for '%s', but got error: %v", tt.ipStr, err)
			}
			if !tt.isValid && err == nil {
				t.Errorf("Expected invalid IP for '%s', but got no error", tt.ipStr)
			}
		})
	}
}

func TestNewServiceInvalidPath(t *testing.T) {
	_, err := NewService("nonexistent/path/to/database.mmdb")
	if err == nil {
		t.Error("Expected error when opening non-existent database, got nil")
	}
}
