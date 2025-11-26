package geoip

import (
	"net"

	"github.com/oschwald/geoip2-golang"
)

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

type LookupData struct {
	Country  string
	City     string
	ISOCode  string
	Timezone string
}

func (s *Service) LookupIP(ipStr string) (*LookupData, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return nil, nil
	}

	record, err := s.DB.City(ip)
	if err != nil {
		return nil, err
	}

	data := &LookupData{
		Country:  record.Country.Names["en"],
		City:     record.City.Names["en"],
		ISOCode:  record.Country.IsoCode,
		Timezone: record.Location.TimeZone,
	}

	return data, nil
}
