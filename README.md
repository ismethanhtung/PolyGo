# PolyGo

High-performance Go server for Polymarket API with caching, WebSocket support, and Swagger UI.

## Features

- **Ultra-fast performance**: Built with Fiber (fasthttp), sonic JSON, and Ristretto cache
- **Complete API coverage**: Markets, Events, Prices, Order Books, Trading, User Data
- **Real-time WebSocket**: Proxy WebSocket connections to Polymarket
- **Smart caching**: TTL-based caching optimized for different data types
- **Swagger UI**: Interactive API documentation
- **Production-ready**: Health checks, rate limiting, graceful shutdown

## Performance Optimizations

| Optimization | Impact |
|-------------|--------|
| Fiber/fasthttp | 10x faster than net/http |
| sonic JSON | 5-10x faster than encoding/json |
| Ristretto cache | High-performance concurrent cache |
| Connection pooling | Reuse HTTP connections |
| Zero-copy responses | Minimal memory allocations |
| Prefork mode | Multi-process for multi-core CPUs |

## Quick Start

### Prerequisites

- Go 1.22+
- Make (optional)

### Install & Run

```bash
# Clone
git clone https://github.com/yourusername/polygo.git
cd polygo

# Install dependencies
go mod download

# Run
go run ./cmd/server

# Or use make
make run
```

Server starts at `http://localhost:8080`

### Swagger UI

Open `http://localhost:8080/swagger/index.html` for interactive API docs.

## API Endpoints

### Public Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/api/v1/markets` | List markets |
| GET | `/api/v1/markets/:id` | Get market by ID |
| GET | `/api/v1/events` | List events |
| GET | `/api/v1/events/:id` | Get event by ID |
| GET | `/api/v1/price/:token_id` | Get current price |
| GET | `/api/v1/book/:token_id` | Get order book |
| GET | `/api/v1/spread/:token_id` | Get spread |
| GET | `/api/v1/top-movers` | Top moving markets |
| GET | `/api/v1/leaderboard` | Trading leaderboard |

### Authenticated Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/v1/orders` | Create order |
| GET | `/api/v1/orders` | List orders |
| DELETE | `/api/v1/orders/:id` | Cancel order |

### WebSocket

PolyGo cung cấp WebSocket endpoints để nhận dữ liệu real-time từ Polymarket. Server tự động kết nối với Polymarket WebSocket và proxy dữ liệu đến clients.

| Endpoint | Description |
|----------|-------------|
| `/ws/market/:market_id` | Subscribe to updates cho một market cụ thể |
| `/ws/markets` | Subscribe to updates cho tất cả markets |

#### WebSocket Usage

**1. Single Market Subscription:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/market/0x1234...');

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Market update:', data);
};
```

**2. All Markets Subscription:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/markets');

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Market update:', data);
};

// Gửi ping để kiểm tra kết nối
ws.send(JSON.stringify({ type: 'ping' }));
```

**3. Dynamic Subscription (Single Market mode):**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/market/0x1234...');

// Subscribe thêm markets
ws.send(JSON.stringify({
    type: 'subscribe',
    markets: ['0x5678...', '0x9abc...']
}));

// Unsubscribe markets
ws.send(JSON.stringify({
    type: 'unsubscribe',
    markets: ['0x5678...']
}));
```

#### Testing WebSocket

Mở file `websocket-test.html` trong trình duyệt để test WebSocket và xem streaming data:

```bash
# Mở file trong trình duyệt
open websocket-test.html
# hoặc
xdg-open websocket-test.html
```

File này cung cấp:
- Giao diện để kết nối WebSocket
- Hiển thị real-time messages
- Xem dữ liệu JSON mới nhất
- Thống kê số lượng messages nhận được

## Configuration

### Environment Variables

```bash
# Server
POLYGO_HOST=0.0.0.0
POLYGO_PORT=8080
POLYGO_DEBUG=false
POLYGO_PREFORK=false

# Polymarket API URLs (defaults provided)
POLYGO_CLOB_URL=https://clob.polymarket.com
POLYGO_GAMMA_URL=https://gamma-api.polymarket.com
POLYGO_DATA_URL=https://data-api.polymarket.com

# Cache
POLYGO_CACHE_MAX_COST=1073741824  # 1GB
POLYGO_CACHE_MARKETS_TTL=30s
POLYGO_CACHE_PRICES_TTL=100ms
```

### Config File

Create `config.yaml`:

```yaml
server:
  host: 0.0.0.0
  port: 8080
  prefork: false
  debug: true

cache:
  max_cost: 1073741824
  markets_ttl: 30s
  events_ttl: 30s
  prices_ttl: 100ms
  order_book_ttl: 50ms
```

## Authentication

For trading endpoints, include these headers:

```
POLY-API-KEY: your-api-key
POLY-API-SECRET: your-api-secret
POLY-PASSPHRASE: your-passphrase
POLY-SIGNATURE: request-signature
POLY-TIMESTAMP: unix-timestamp
```

## Development

### Commands

```bash
make build          # Build binary
make run            # Run server
make test           # Run all tests
make test-unit      # Run unit tests
make bench          # Run benchmarks
make lint           # Run linter
make swagger        # Generate Swagger docs
make docker-build   # Build Docker image
```

### Project Structure

```
polygo/
├── cmd/server/          # Entry point
├── internal/
│   ├── api/
│   │   ├── handlers/    # HTTP handlers
│   │   ├── middleware/  # Middleware
│   │   └── routes.go    # Route definitions
│   ├── polymarket/      # Polymarket clients
│   ├── cache/           # Cache layer
│   ├── config/          # Configuration
│   └── models/          # Data models
├── pkg/response/        # Response utilities
├── docs/                # Swagger docs
├── tests/               # Tests
├── Dockerfile
├── Makefile
└── README.md
```

## Docker

```bash
# Build image
docker build -t polygo .

# Run container
docker run -p 8080:8080 polygo

# Or use docker-compose
docker-compose up -d
```

## Benchmarks

Run benchmarks:

```bash
make bench
```

Expected performance:
- Health check: ~50μs
- Cached response: ~100μs
- API proxy: ~5-50ms (depends on Polymarket)

## License

MIT License
