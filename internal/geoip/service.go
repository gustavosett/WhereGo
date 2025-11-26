package geoip

import (
	"errors"
	"net/netip"

	"github.com/oschwald/geoip2-golang"
)

var ErrInvalidIP = errors.New("invalid IP address")

type Service struct {
	DB *geoip2.Reader
}

func NewService(dbPath string) (*Service, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Service{DB: db}, nil
}

type LookupResult struct {
	IP       string `json:"ip"`
	Country  string `json:"country"`
	City     string `json:"city"`
	ISOCode  string `json:"iso_code"`
	Timezone string `json:"timezone"`
}

func (s *Service) LookupIP(ipStr string) (*LookupResult, error) {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, ErrInvalidIP
	}

	record, err := s.DB.City(addr.AsSlice())
	if err != nil {
		return nil, err
	}

	return &LookupResult{
		IP:       ipStr,
		Country:  record.Country.Names["en"],
		City:     record.City.Names["en"],
		ISOCode:  record.Country.IsoCode,
		Timezone: record.Location.TimeZone,
	}, nil
}
