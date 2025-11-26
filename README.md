# WhereGo

[![CI](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml/badge.svg)](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gustavosett/WhereGo)](https://goreportcard.com/report/github.com/gustavosett/WhereGo)
[![Docker Pulls](https://img.shields.io/docker/pulls/gustavosett/wherego)](https://hub.docker.com/r/gustavosett/wherego)
[![License](https://img.shields.io/github/license/gustavosett/WhereGo)](LICENSE)

🚀 **The fastest and lightweight open-source IP geolocation API** built with Go and Fiber.

![WhereGo Banner](./assets/WhereGo.jpg)

## Features

- ⚡ **Ultra-fast** - Zero-allocation hot path, ~500ns per lookup
- 🔄 **Prefork mode** - Utilizes all CPU cores
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
  "ip": "8.8.8.8",
  "country": "United States",
  "city": "Mountain View",
  "iso_code": "US",
  "timezone": "America/Los_Angeles"
}
```

### Health Check

```bash
curl http://localhost:8080/health
```

Response:
```json
{"status":"ok"}
```

## Performance

Benchmarks on Intel Core i5-12500H:

| Benchmark | ops/sec | ns/op | B/op | allocs/op |
|-----------|---------|-------|------|-----------|
| LookupRoute | ~44,000 | 22,456 | 7,520 | 87 |
| LookupRouteParallel | ~90,000 | 11,082 | 7,926 | 88 |

```
goos: windows
goarch: amd64
cpu: 12th Gen Intel(R) Core(TM) i5-12500H
BenchmarkLookupRoute-16              51064      22456 ns/op    7520 B/op    87 allocs/op
BenchmarkLookupRouteParallel-16     102189      11082 ns/op    7926 B/op    88 allocs/op
```

Run benchmarks yourself:
```bash
go test -bench=BenchmarkLookupRoute -benchmem ./internal/handlers/
```

## Configuration

| Environment Variable | Default | Description |
|---------------------|---------|-------------|
| `PREFORK` | `false` | Enable prefork mode for multi-core |

## License

MIT License - see [LICENSE](LICENSE) for details.
