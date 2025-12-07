package geoip

import (
	"errors"
	"net/netip"
)

var ErrInvalidIP = errors.New("invalid IP address")

type Service struct {
	DB *Reader
}

func NewService(dbPath string) (*Service, error) {
	db, err := Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Service{DB: db}, nil
}

func (s *Service) LookupIP(ipStr string) (*City, error) {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return nil, ErrInvalidIP
	}

	return s.DB.City(addr)
}
