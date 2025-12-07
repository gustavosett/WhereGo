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
    "traits": {
        "ip_address": "8.8.8.8",
        "network": "8.8.8.0/24",
        "is_anycast": false
    },
    "postal": {
        "code": ""
    },
    "continent": {
        "names": {
            "de": "Nordamerika",
            "en": "North America",
            "es": "Norteamérica",
            "fr": "Amérique du Nord",
            "ja": "北アメリカ",
            "pt-BR": "América do Norte",
            "ru": "Северная Америка",
            "zh-CN": "北美洲"
        },
        "code": "NA",
        "geoname_id": 6255149
    },
    "city": {
        "names": {
            "de": "",
            ...
        },
        "geoname_id": 0
    },
    "subdivisions": null,
    "represented_country": {
        "names": {
            "de": "",
            ...
        },
        "iso_code": "",
        "type": "",
        "geoname_id": 0,
        "is_in_european_union": false
    },
    "country": {
        "names": {
            "de": "USA",
            "en": "United States",
            "es": "Estados Unidos",
            "fr": "États Unis",
            "ja": "アメリカ",
            "pt-BR": "EUA",
            "ru": "США",
            "zh-CN": "美国"
        },
        "iso_code": "US",
        "geoname_id": 6252001,
        "is_in_european_union": false
    },
    "registered_country": {
        "names": {
            "de": "USA",
            "en": "United States",
            "es": "Estados Unidos",
            "fr": "États Unis",
            "ja": "アメリカ",
            "pt-BR": "EUA",
            "ru": "США",
            "zh-CN": "美国"
        },
        "iso_code": "US",
        "geoname_id": 6252001,
        "is_in_european_union": false
    },
    "location": {
        "latitude": 37.751,
        "longitude": -97.822,
        "time_zone": "America/Chicago",
        "metro_code": 0,
        "accuracy_radius": 1000
    }
}
```

## Performance

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
