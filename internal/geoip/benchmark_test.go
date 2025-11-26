package geoip

import (
	"testing"
)

const (
	testDBPath   = "../../data/city.db"
	skipMessage  = "Database not available for benchmarking"
	testIPGoogle = "8.8.8.8"
)

func BenchmarkLookupIP(b *testing.B) {
	service, err := NewService(testDBPath)
	if err != nil {
		b.Skip(skipMessage)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	data := service.GetLookupData()
	defer service.PutLookupData(data)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = service.LookupIP(testIPGoogle, data)
	}
}

func BenchmarkLookupIPParallel(b *testing.B) {
	service, err := NewService(testDBPath)
	if err != nil {
		b.Skip(skipMessage)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		data := service.GetLookupData()
		defer service.PutLookupData(data)

		for pb.Next() {
			_ = service.LookupIP(testIPGoogle, data)
		}
	})
}

func BenchmarkParseIP(b *testing.B) {
	service, err := NewService(testDBPath)
	if err != nil {
		b.Skip(skipMessage)
		return
	}
	defer service.DB.Close() //nolint:errcheck

	data := service.GetLookupData()
	defer service.PutLookupData(data)

	ips := []string{
		testIPGoogle,
		"1.1.1.1",
		"208.67.222.222",
		"9.9.9.9",
		"185.228.168.9",
	}

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = service.LookupIP(ips[i%len(ips)], data)
	}
}
