# WhereGo

[![CI](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml/badge.svg)](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gustavosett/WhereGo)](https://goreportcard.com/report/github.com/gustavosett/WhereGo)
[![Docker Pulls](https://img.shields.io/docker/pulls/gustavosett/wherego)](https://hub.docker.com/r/gustavosett/wherego)
[![License](https://img.shields.io/github/license/gustavosett/WhereGo)](LICENSE)

🚀 **The fastest and lightweight open-source IP geolocation API** built with Go and Echo.

![WhereGo Banner](./assets/WhereGo.jpg)

## Features

- ⚡ **Ultra-fast** - ~1ms per lookup
- 🌐 **Multi-language** - Location names in 8 languages
- 🐳 **Docker ready** - Multi-arch images (amd64/arm64)
- 📦 **Tiny image** - ~40MB distroless container
- 🛡️ **Production ready** - Health checks, graceful shutdown

## Quick Start

### Docker

```bash
docker run -p 8080:8080 gustavosett/wherego:latest
```

### From Source

```bash
git clone https://github.com/gustavosett/WhereGo.git
cd WhereGo
go run ./cmd/api
```

## API Endpoints

### Lookup IP

```bash
curl http://localhost:8080/lookup/8.8.8.8
```

Response:
```json
{
  "City": {
    "Names": {"en": "Mountain View", ...},
    "GeoNameID": 5375480
  },
  "Country": {
    "Names": {"en": "United States", "es": "Estados Unidos", ...},
    "IsoCode": "US",
    "GeoNameID": 6252001,
    "IsInEuropeanUnion": false
  },
  "Continent": {
    "Names": {"en": "North America", ...},
    "Code": "NA",
    "GeoNameID": 6255149
  },
  "Location": {
    "TimeZone": "America/Los_Angeles",
    "Latitude": 37.386,
    "Longitude": -122.0838,
    "MetroCode": 807,
    "AccuracyRadius": 1000
  },
  "Postal": {"Code": "94035"},
  "Subdivisions": [{"Names": {"en": "California"}, "IsoCode": "CA", ...}],
  "RegisteredCountry": {...},
  "RepresentedCountry": {...},
  "Traits": {
    "IsAnonymousProxy": false,
    "IsAnycast": false,
    "IsSatelliteProvider": false
  }
}
```

## Performance

![WhereGo Benchmark](./assets/benchmark.jpg)

Benchmarks on Intel Core i5-12500H:

| Benchmark | ops/sec | ns/op | B/op | allocs/op |
|-----------|---------|-------|------|-----------|
| Lookup | ~107,000 | 11,093 | 8,131 | 85 |
| LookupParallel | ~260,000 | 4,364 | 8,129 | 85 |

```
goos: windows
goarch: amd64
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkLookup-16                107030             11093 ns/op            8131 B/op         85 allocs/op
BenchmarkLookupParallel-16        260906              4364 ns/op            8129 B/op         85 allocs/op
```

Run benchmarks yourself:
```bash
go test -bench=Benchmark -benchmem -run=^$ ./internal/handlers/
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PORT` | `8080` | Server port |

## License

MIT License - see [LICENSE](LICENSE) for details.
