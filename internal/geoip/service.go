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

func (s *Service) LookupIP(ipStr string) (*geoip2.City, error) {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, ErrInvalidIP
	}

	return s.DB.City(addr.AsSlice())
}
