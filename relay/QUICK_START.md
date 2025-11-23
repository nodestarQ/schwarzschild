# Quick Start Guide

## 1. Setup

```bash
# Copy environment file
cp .env.example .env

# Edit .env with your values:
# - ETHEREUM_RPC: Your Sepolia RPC endpoint (from Infura, Alchemy, etc.)
# - CONTRACT_ADDRESS: Your deployed contract address
# - PRIVATE_KEY: Your wallet private key (without 0x prefix)
# - CONTRACT_ABI: Your smart contract's ABI (as JSON string)
```

## 2. Run with Docker Compose (Recommended)

```bash
# Start the service with hot reload
docker-compose up

# Service will be available at http://localhost:8080
```

## 3. Test the API

```bash
# Health check
curl http://localhost:8080/health

# Submit a relay transaction (StealthBurnRegistry.emitBurn)
curl -X POST http://localhost:8080/relay \
  -H "Content-Type: application/json" \
  -d '{
    "ephemeralPublicKey": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
    "burnAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f42eD5"
  }'
```

**Parameter Format:**
- `ephemeralPublicKey`: 64-character hex string (32 bytes) - with or without 0x prefix
- `burnAddress`: Ethereum address format (0x + 40 hex characters)

## 4. Local Development (Without Docker)

```bash
# Install air for hot reload
go install github.com/cosmtrek/air@latest

# Run with auto-reload
air -c .air.toml

# Or run directly
go run ./cmd/main.go
```

## 5. Production Deployment

```bash
# Build Docker image
docker build -t relay:latest .

# Run container
docker run -p 8080:8080 --env-file .env relay:latest
```

## Getting Testnet Credentials

### Sepolia RPC Endpoint
1. Visit https://infura.io
2. Create account and project
3. Select Sepolia network
4. Copy API key: `https://sepolia.infura.io/v3/YOUR_KEY`

### Sepolia Testnet ETH
1. Visit https://sepoliafaucet.com
2. Enter your wallet address
3. Claim free testnet ETH

### Private Key
Export from MetaMask:
1. Settings → Security & Privacy
2. Click "Reveal Private Key"
3. Copy without the "0x" prefix
4. Add to `.env` as `PRIVATE_KEY`

## API Response Examples

### Success
```json
{
  "success": true,
  "message": "Transaction submitted successfully",
  "txHash": "0x1234567890abcdef..."
}
```

### Error
```json
{
  "success": false,
  "error": "ephemeralPublicKey and burnAddress are required"
}
```

## About the Contract Call

This relay calls the **StealthBurnRegistry** contract's `emitBurn` function:

```solidity
function emitBurn(bytes32 ephemeralPublicKey, address burnAddress) public
```

**Automatic Conversion:**
- `ephemeralPublicKey` string input is automatically converted to bytes32
- `burnAddress` string is validated and converted to Ethereum address
- Supports both with and without "0x" prefix for hex values

If you need to customize for a different contract:

1. Update `CONTRACT_ADDRESS` in `.env`
2. Update `CONTRACT_ABI` in `.env` with your contract's full ABI
3. If the method name differs, edit `internal/ethereum/client.go` to change "emitBurn" on line ~90

## Troubleshooting

| Issue | Solution |
|-------|----------|
| Connection refused | Make sure service is running on port 8080 |
| "CONTRACT_ADDRESS is required" | Add CONTRACT_ADDRESS to .env |
| "Failed to connect to Ethereum" | Check ETHEREUM_RPC is correct and accessible |
| "Failed to get nonce" | Account needs to exist on Sepolia (has no balance) |
| "Failed to pack contract call" | Check CONTRACT_ABI is valid JSON |

For more details, see [README.md](README.md)
