package geoip

import (
	"net"
	"testing"
)

func TestLookupData(t *testing.T) {
	// Test that LookupData struct can be created properly
	data := &LookupData{
		Country:  "United States",
		City:     "New York",
		ISOCode:  "US",
		Timezone: "America/New_York",
	}

	if data.Country != "United States" {
		t.Errorf("Expected Country 'United States', got '%s'", data.Country)
	}
	if data.City != "New York" {
		t.Errorf("Expected City 'New York', got '%s'", data.City)
	}
	if data.ISOCode != "US" {
		t.Errorf("Expected ISOCode 'US', got '%s'", data.ISOCode)
	}
	if data.Timezone != "America/New_York" {
		t.Errorf("Expected Timezone 'America/New_York', got '%s'", data.Timezone)
	}
}

func TestLookupDataReset(t *testing.T) {
	data := &LookupData{
		Country:  "United States",
		City:     "New York",
		ISOCode:  "US",
		Timezone: "America/New_York",
	}

	data.Reset()

	if data.Country != "" {
		t.Errorf("Expected empty Country after Reset, got '%s'", data.Country)
	}
	if data.City != "" {
		t.Errorf("Expected empty City after Reset, got '%s'", data.City)
	}
	if data.ISOCode != "" {
		t.Errorf("Expected empty ISOCode after Reset, got '%s'", data.ISOCode)
	}
	if data.Timezone != "" {
		t.Errorf("Expected empty Timezone after Reset, got '%s'", data.Timezone)
	}
}

func TestErrInvalidIP(t *testing.T) {
	if ErrInvalidIP == nil {
		t.Error("ErrInvalidIP should not be nil")
	}
	if ErrInvalidIP.Error() != "invalid IP address" {
		t.Errorf("Expected 'invalid IP address', got '%s'", ErrInvalidIP.Error())
	}
}

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
			ip := net.ParseIP(tt.ipStr)
			if tt.isValid && ip == nil {
				t.Errorf("Expected valid IP for '%s', but got nil", tt.ipStr)
			}
			if !tt.isValid && ip != nil {
				t.Errorf("Expected invalid IP for '%s', but got %v", tt.ipStr, ip)
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

// MockService provides a way to test without a real database
type MockLookupData struct {
	Country  string
	City     string
	ISOCode  string
	Timezone string
}

func TestLookupIPInvalidIP(t *testing.T) {
	// Since we can't easily mock the geoip2 database,
	// we test the IP parsing logic separately
	invalidIPs := []string{
		"",
		"invalid",
		"256.256.256.256",
		"abc.def.ghi.jkl",
		"192.168.1",
	}

	for _, ipStr := range invalidIPs {
		ip := net.ParseIP(ipStr)
		if ip != nil {
			t.Errorf("Expected nil for invalid IP '%s', got %v", ipStr, ip)
		}
	}
}

func TestLookupIPValidIPFormat(t *testing.T) {
	validIPs := []string{
		"8.8.8.8",
		"1.1.1.1",
		"192.168.1.1",
		"10.0.0.1",
		"172.16.0.1",
		"2001:4860:4860::8888",
		"::1",
	}

	for _, ipStr := range validIPs {
		ip := net.ParseIP(ipStr)
		if ip == nil {
			t.Errorf("Expected valid IP for '%s', got nil", ipStr)
		}
	}
}
