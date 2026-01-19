# WebSocket Documentation

## Tổng quan

PolyGo cung cấp WebSocket endpoints để nhận dữ liệu real-time từ Polymarket. Server tự động kết nối với Polymarket WebSocket (CLOB và Live Data) và proxy dữ liệu đến clients của bạn.

## Cách hoạt động

1. **Server kết nối với Polymarket**: Khi server khởi động, nó tự động kết nối với 2 WebSocket endpoints của Polymarket:
   - CLOB WebSocket: `wss://ws-subscriptions-clob.polymarket.com/ws/`
   - Live Data WebSocket: `wss://ws-live-data.polymarket.com`

2. **Client kết nối với PolyGo**: Client của bạn kết nối với PolyGo WebSocket endpoints

3. **Proxy dữ liệu**: PolyGo nhận dữ liệu từ Polymarket và forward đến tất cả clients đã subscribe

## WebSocket Endpoints

### 1. Single Market Subscription

**Endpoint:** `ws://localhost:8080/ws/market/:market_id`

**Mô tả:** Subscribe để nhận updates cho một market cụ thể

**Ví dụ:**
```javascript
const marketId = '0x1234567890abcdef...';
const ws = new WebSocket(`ws://localhost:8080/ws/market/${marketId}`);

ws.onopen = () => {
    console.log('Đã kết nối!');
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Market update:', data);
    // Dữ liệu sẽ tự động update khi có thay đổi từ Polymarket
};

ws.onerror = (error) => {
    console.error('WebSocket error:', error);
};

ws.onclose = () => {
    console.log('Đã ngắt kết nối');
};
```

**Subscribe thêm markets (dynamic subscription):**
```javascript
// Sau khi đã kết nối, bạn có thể subscribe thêm markets
ws.send(JSON.stringify({
    type: 'subscribe',
    markets: ['0x5678...', '0x9abc...']
}));
```

**Unsubscribe markets:**
```javascript
ws.send(JSON.stringify({
    type: 'unsubscribe',
    markets: ['0x5678...']
}));
```

**Ping/Pong:**
```javascript
// Gửi ping để kiểm tra kết nối
ws.send(JSON.stringify({
    type: 'ping'
}));

// Server sẽ trả về pong với timestamp
// {
//   "type": "pong",
//   "timestamp": 1234567890123
// }
```

### 2. All Markets Subscription

**Endpoint:** `ws://localhost:8080/ws/markets`

**Mô tả:** Subscribe để nhận updates cho TẤT CẢ markets

**Ví dụ:**
```javascript
const ws = new WebSocket('ws://localhost:8080/ws/markets');

ws.onopen = () => {
    console.log('Đã kết nối! Bạn sẽ nhận updates từ tất cả markets');
    
    // Gửi ping để kiểm tra
    ws.send(JSON.stringify({ type: 'ping' }));
};

ws.onmessage = (event) => {
    const data = JSON.parse(event.data);
    console.log('Market update từ bất kỳ market nào:', data);
    // Dữ liệu sẽ tự động update khi có thay đổi từ Polymarket
};
```

## Dữ liệu nhận được

Dữ liệu nhận được từ WebSocket sẽ có format tùy thuộc vào loại update từ Polymarket:

### Market Updates
```json
{
  "type": "update",
  "channel": "market",
  "markets": ["0x1234..."],
  "data": {
    "price": "0.65",
    "volume": "1000",
    "liquidity": "5000",
    ...
  }
}
```

### Price Updates
```json
{
  "type": "price",
  "market": "0x1234...",
  "data": {
    "token_id": "0x5678...",
    "price": "0.65",
    "timestamp": 1234567890
  }
}
```

### Order Book Updates
```json
{
  "type": "orderbook",
  "market": "0x1234...",
  "data": {
    "bids": [...],
    "asks": [...],
    "timestamp": 1234567890
  }
}
```

## Testing với HTML Client

File `websocket-test.html` cung cấp một giao diện đầy đủ để test WebSocket:

1. **Mở file trong trình duyệt:**
   ```bash
   open websocket-test.html
   ```

2. **Nhập thông tin:**
   - Chọn chế độ: Single Market hoặc All Markets
   - Nhập Market ID (nếu chọn Single Market)
   - Nhập Server URL (mặc định: `ws://localhost:8080`)

3. **Kết nối:**
   - Click "Kết nối"
   - Xem messages real-time trong panel bên trái
   - Xem dữ liệu JSON mới nhất trong panel bên phải
   - Theo dõi thống kê: tổng messages, lỗi, dữ liệu nhận được

## Client Messages

### Subscribe Message
```json
{
  "type": "subscribe",
  "markets": ["0x1234...", "0x5678..."]
}
```

### Unsubscribe Message
```json
{
  "type": "unsubscribe",
  "markets": ["0x1234...", "0x5678..."]
}
```

### Ping Message
```json
{
  "type": "ping"
}
```

## Server Messages

### Pong Response
```json
{
  "type": "pong",
  "timestamp": 1234567890123
}
```

### Market Data
Dữ liệu từ Polymarket được forward nguyên bản đến client. Format phụ thuộc vào loại update.

## Error Handling

WebSocket sẽ tự động reconnect nếu mất kết nối với Polymarket. Client nên implement reconnection logic:

```javascript
let ws;
let reconnectInterval = 1000; // Start with 1 second

function connect() {
    ws = new WebSocket('ws://localhost:8080/ws/markets');
    
    ws.onclose = () => {
        console.log('Mất kết nối, đang reconnect...');
        setTimeout(connect, reconnectInterval);
        reconnectInterval = Math.min(reconnectInterval * 2, 30000); // Max 30s
    };
    
    ws.onerror = (error) => {
        console.error('WebSocket error:', error);
    };
}

connect();
```

## Health Check

Kiểm tra trạng thái WebSocket connection:

```bash
curl http://localhost:8080/health
```

Response sẽ bao gồm:
```json
{
  "status": "ok",
  "services": {
    "websocket": "connected" // hoặc "disconnected"
  }
}
```

## Best Practices

1. **Reconnection**: Luôn implement reconnection logic trong client
2. **Error Handling**: Xử lý errors và connection drops
3. **Rate Limiting**: Không gửi quá nhiều messages trong thời gian ngắn
4. **Cleanup**: Đóng connection khi không cần thiết
5. **Ping/Pong**: Sử dụng ping để kiểm tra kết nối định kỳ

## Troubleshooting

### Không nhận được dữ liệu

1. Kiểm tra server đã kết nối với Polymarket:
   ```bash
   curl http://localhost:8080/health
   ```

2. Kiểm tra Market ID có đúng không:
   ```bash
   curl http://localhost:8080/api/v1/markets/0x1234...
   ```

3. Kiểm tra logs của server để xem có lỗi gì không

### Connection bị đóng

- Server sẽ tự động reconnect với Polymarket
- Client nên implement reconnection logic
- Kiểm tra network connectivity

### Dữ liệu không update

- Đảm bảo market đang active và có trading activity
- Kiểm tra subscription đã đúng market ID chưa
- Xem logs để kiểm tra có messages từ Polymarket không
