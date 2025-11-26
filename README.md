# WhereGo

[![CI](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml/badge.svg)](https://github.com/gustavosett/WhereGo/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/gustavosett/WhereGo)](https://goreportcard.com/report/github.com/gustavosett/WhereGo)
[![Docker Pulls](https://img.shields.io/docker/pulls/gustavosett/wherego)](https://hub.docker.com/r/gustavosett/wherego)
[![License](https://img.shields.io/github/license/gustavosett/WhereGo)](LICENSE)

🚀 **The fastest and lightweight open-source IP geolocation API** built with Go and Echo.

![WhereGo Banner](./assets/WhereGo.jpg)

## Features

- ⚡ **Ultra-fast** - ~1ms per lookup
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
