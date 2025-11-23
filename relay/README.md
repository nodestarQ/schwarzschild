# Ethereum Relay Service

A production-ready Go-based relay service that listens for API requests and relays transactions to the Ethereum Sepolia testnet. The service accepts ephemeral public keys and burn addresses, calls a smart contract with these parameters, and submits the transaction to the blockchain.

## Features

- ✅ RESTful API endpoint for relay transactions
- ✅ Ethereum Sepolia testnet integration
- ✅ Smart contract interaction with customizable ABI
- ✅ Hot reload development with air
- ✅ Docker containerization (development and production)
- ✅ Environment-based configuration
- ✅ Health check endpoint
- ✅ Comprehensive error handling

## Prerequisites

- Go 1.24+ (for local development)
- Docker & Docker Compose (optional, for containerized development)
- Sepolia testnet RPC endpoint (Infura, Alchemy, etc.)
- Private key with Sepolia testnet funds (for gas fees)
- Smart contract deployed on Sepolia testnet

## Quick Start

### Using Docker Compose (Recommended for Development)

1. Clone the repository and navigate to the project directory:
```bash
cd relay
```

2. Copy the environment configuration:
```bash
cp .env.example .env
```

3. Edit `.env` and add your configuration:
```bash
PORT=8080
ETHEREUM_RPC=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
CONTRACT_ADDRESS=0x...your_contract_address...
PRIVATE_KEY=your_private_key_without_0x_prefix
CONTRACT_ABI=[{"inputs":[{"name":"ephemeralPublicKey","type":"string"},{"name":"burnAddress","type":"address"}],"name":"relay","outputs":[],"stateMutability":"nonpayable","type":"function"}]
```

4. Start the service with hot reload:
```bash
docker-compose up
```

The server will be running on `http://localhost:8080`

### Local Development (Without Docker)

1. Set up environment:
```bash
cp .env.example .env
# Edit .env with your configuration
```

2. Install dependencies:
```bash
go mod download
```

3. Run with air for hot reload:
```bash
go install github.com/cosmtrek/air@latest
air -c .air.toml
```

4. Or run directly:
```bash
go run ./cmd/main.go
```

### Production Deployment

1. Build the Docker image:
```bash
docker build -t relay:latest .
```

2. Run the container:
```bash
docker run -p 8080:8080 --env-file .env relay:latest
```

## API Endpoints

### Relay Transaction

Submit a transaction to be relayed to the smart contract.

**Request:**
```bash
curl -X POST http://localhost:8080/relay \
  -H "Content-Type: application/json" \
  -d '{
    "ephemeralPublicKey": "0x1234567890abcdef...",
    "burnAddress": "0xabcdefg1234567890..."
  }'
```

**Response (Success):**
```json
{
  "success": true,
  "message": "Transaction submitted successfully",
  "txHash": "0x1234567890abcdef..."
}
```

**Response (Error):**
```json
{
  "success": false,
  "error": "Failed to relay transaction: <error_message>"
}
```

### Health Check

Check if the service is running.

**Request:**
```bash
curl http://localhost:8080/health
```

**Response:**
```json
{
  "status": "ok"
}
```

## Configuration

All configuration is managed through environment variables in the `.env` file:

| Variable | Description | Example |
|----------|-------------|---------|
| `PORT` | Server port | `8080` |
| `ETHEREUM_RPC` | Sepolia RPC endpoint | `https://sepolia.infura.io/v3/YOUR_KEY` |
| `CONTRACT_ADDRESS` | Smart contract address | `0x...` |
| `PRIVATE_KEY` | Wallet private key (no 0x prefix) | `abc123def456...` |
| `CONTRACT_ABI` | Smart contract ABI (JSON) | `[{...}]` |

## Getting Sepolia Testnet Credentials

### 1. Get an Infura API Key

1. Visit [infura.io](https://infura.io)
2. Sign up and create a new project
3. Select Sepolia as your network
4. Copy your API key

Your RPC endpoint will be:
```
https://sepolia.infura.io/v3/YOUR_API_KEY
```

### 2. Get Sepolia Testnet Funds

1. Visit [Sepolia Faucet](https://sepoliafaucet.com)
2. Enter your wallet address
3. Receive free testnet ETH

### 3. Export Your Private Key

From your wallet (MetaMask, etc.):
1. Open your wallet settings
2. Find "Export Private Key"
3. Copy without the "0x" prefix
4. Add to `.env` as `PRIVATE_KEY`

## Smart Contract Considerations

The relay service calls a contract method named `relay` with the following signature:

```solidity
function relay(string memory ephemeralPublicKey, address burnAddress) public
```

If your contract uses a different method name or parameters, you'll need to:

1. Update the method name in `internal/ethereum/client.go` (line ~86):
```go
data, err := c.contractABI.Pack("yourMethodName", param1, param2)
```

2. Update the contract ABI in your `.env` file to match your actual contract

## Development

### Project Structure

```
relay/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go           # Configuration management
│   ├── ethereum/
│   │   └── client.go           # Ethereum client and contract interaction
│   └── handlers/
│       └── relayer.go          # HTTP request handlers
├── go.mod                       # Go module definition
├── go.sum                       # Go module checksums
├── Dockerfile                   # Production Docker image
├── Dockerfile.dev              # Development Docker image with air
├── docker-compose.yml          # Docker Compose configuration
├── .air.toml                   # Air hot reload configuration
├── .env.example                # Example environment variables
└── README.md                   # This file
```

### Hot Reload Development

The project is configured to automatically reload when you modify Go files:

1. Using Docker Compose:
```bash
docker-compose up
```

2. Using air directly:
```bash
air -c .air.toml
```

### Building

Build the application:
```bash
go build ./cmd/main.go
```

### Testing

To test the relay endpoint locally:

```bash
# Make sure the server is running
go run ./cmd/main.go

# In another terminal, send a test request:
curl -X POST http://localhost:8080/relay \
  -H "Content-Type: application/json" \
  -d '{
    "ephemeralPublicKey": "test_public_key",
    "burnAddress": "0x0000000000000000000000000000000000000000"
  }'
```

## Troubleshooting

### "CONTRACT_ADDRESS is required"
- Make sure you've set `CONTRACT_ADDRESS` in your `.env` file
- Ensure the address is valid and starts with `0x`

### "PRIVATE_KEY is required"
- Set `PRIVATE_KEY` in your `.env` file
- Don't include the `0x` prefix
- Make sure the account has Sepolia testnet ETH for gas fees

### "Failed to connect to Ethereum"
- Verify your `ETHEREUM_RPC` endpoint is correct
- Check that you have internet connectivity
- For Infura, ensure your API key is valid

### "Failed to get nonce"
- The account may not exist on Sepolia yet
- Send a test transaction to create it first
- Ensure the account has funds

### "Failed to pack contract call"
- Check that your `CONTRACT_ABI` is valid JSON
- Verify the ABI matches your actual contract
- Ensure the contract has a `relay` method (or update the code for your method name)

## Environment Setup Example

Here's a complete `.env` example with placeholder values:

```
PORT=8080
ETHEREUM_RPC=https://sepolia.infura.io/v3/abc123def456ghi789
CONTRACT_ADDRESS=0x1234567890123456789012345678901234567890
PRIVATE_KEY=1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef
CONTRACT_ABI=[{"inputs":[{"internalType":"string","name":"ephemeralPublicKey","type":"string"},{"internalType":"address","name":"burnAddress","type":"address"}],"name":"relay","outputs":[],"stateMutability":"nonpayable","type":"function"}]
```

## Deployment Notes

### Security Considerations

- Never commit `.env` file to version control
- Use secure environment variable management in production
- Rotate private keys regularly
- Use a dedicated account for relay operations
- Monitor transaction history and gas costs

### Performance

- Default gas limit is 300,000. Adjust based on your contract
- Consider implementing gas price optimization for production
- Monitor Sepolia network congestion

### Monitoring

- Check logs for transaction failures
- Monitor gas usage and costs
- Set up alerts for failed transactions
- Track contract interaction patterns

## License

MIT

## Support

For issues or questions, please open an issue on the repository.
