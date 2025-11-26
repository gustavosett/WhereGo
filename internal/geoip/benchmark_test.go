package geoip

import (
	"testing"
)

func BenchmarkLookupIP(b *testing.B) {
	service, err := NewService("../../data/city.db")
	if err != nil {
		b.Skip("Database not available for benchmarking")
		return
	}
	defer service.DB.Close()

	data := service.GetLookupData()
	defer service.PutLookupData(data)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = service.LookupIP("8.8.8.8", data)
	}
}

func BenchmarkLookupIPParallel(b *testing.B) {
	service, err := NewService("../../data/city.db")
	if err != nil {
		b.Skip("Database not available for benchmarking")
		return
	}
	defer service.DB.Close()

	b.ResetTimer()
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		data := service.GetLookupData()
		defer service.PutLookupData(data)

		for pb.Next() {
			_ = service.LookupIP("8.8.8.8", data)
		}
	})
}

func BenchmarkParseIP(b *testing.B) {
	service, err := NewService("../../data/city.db")
	if err != nil {
		b.Skip("Database not available for benchmarking")
		return
	}
	defer service.DB.Close()

	data := service.GetLookupData()
	defer service.PutLookupData(data)

	ips := []string{
		"8.8.8.8",
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
