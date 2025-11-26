package geoip

import (
	"errors"
	"net/netip"
	"sync"

	"github.com/oschwald/geoip2-golang"
)

var ErrInvalidIP = errors.New("invalid IP address")

type Service struct {
	DB   *geoip2.Reader
	pool sync.Pool
}

func NewService(dbPath string) (*Service, error) {
	db, err := geoip2.Open(dbPath)
	if err != nil {
		return nil, err
	}
	return &Service{
		DB: db,
		pool: sync.Pool{
			New: func() any {
				return &LookupData{}
			},
		},
	}, nil
}

type LookupData struct {
	Country  string
	City     string
	ISOCode  string
	Timezone string
}

// GetLookupData gets a LookupData from the pool
func (s *Service) GetLookupData() *LookupData {
	return s.pool.Get().(*LookupData)
}

// PutLookupData returns a LookupData to the pool
func (s *Service) PutLookupData(d *LookupData) {
	s.pool.Put(d)
}

func (s *Service) LookupIP(ipStr string, data *LookupData) error {
	addr, err := netip.ParseAddr(ipStr)
	if err != nil {
		return ErrInvalidIP
	}

	record, err := s.DB.City(addr.AsSlice())
	if err != nil {
		return err
	}

	data.Country = record.Country.Names["en"]
	data.City = record.City.Names["en"]
	data.ISOCode = record.Country.IsoCode
	data.Timezone = record.Location.TimeZone

	return nil
}
