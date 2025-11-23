# Ethereum Contract Integration - StealthBurnRegistry

## Changes Made

### 1. Updated Contract Function
**File:** `internal/ethereum/client.go`

Changed from generic "relay" function to **StealthBurnRegistry** `emitBurn` function:

```go
// Old (generic):
data, err := c.contractABI.Pack("relay", ephemeralPublicKey, burnAddr)

// New (StealthBurnRegistry):
data, err := c.contractABI.Pack("emitBurn", ephemeralBytes, burnAddr)
```

### 2. Type Conversion for bytes32
Added helper functions to convert the string parameter to proper `bytes32` type:

```go
func parseBytes32(hexStr string) [32]byte
func min(a, b int) int
```

**Capabilities:**
- Accepts hex strings with or without "0x" prefix
- Automatically pads short values with leading zeros
- Converts to proper [32]byte array for Solidity compatibility

### 3. Updated Contract ABI
**File:** `.env.example`

Now includes the complete StealthBurnRegistry ABI:
- Constructor
- StealthBurn event definition
- emitBurn function signature

### 4. Improved Logging
Added detailed logging for debugging:
```go
log.Printf("Packed data for emitBurn: ephemeralPublicKey=%x, burnAddress=%s", ephemeralBytes, burnAddr.Hex())
```

## API Contract

### Request Format
```json
{
  "ephemeralPublicKey": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
  "burnAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f42eD5"
}
```

### Parameter Specifications

**ephemeralPublicKey**
- Type: bytes32 (32-byte hex value)
- Input: String (hex with or without 0x prefix)
- Example: `0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef`
- Can also be: `1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef`
- Shorter values are padded with leading zeros

**burnAddress**
- Type: address (20-byte Ethereum address)
- Input: String (hex with 0x prefix)
- Example: `0x742d35Cc6634C0532925a3b844Bc9e7595f42eD5`
- Must be a valid Ethereum address format

### Response Format
```json
{
  "success": true,
  "message": "Transaction submitted successfully",
  "txHash": "0x123abc..."
}
```

## Smart Contract Details

**Contract Name:** StealthBurnRegistry
**Network:** Ethereum Sepolia Testnet
**Chain ID:** 11155111

### Function Signature
```solidity
function emitBurn(
  bytes32 ephemeralPublicKey,
  address burnAddress
) public nonpayable
```

### Event Emitted
```solidity
event StealthBurn(
  bytes32 indexed ephemeralPublicKey,
  address indexed burnAddress
)
```

## Gas Configuration

**Gas Limit:** 100,000
- Adjusted based on the emitBurn function requirements
- Can be modified in `internal/ethereum/client.go` line 83

**Gas Price:** Dynamic
- Automatically retrieved from the network using `SuggestGasPrice()`

## Testing

### Example cURL Request
```bash
curl -X POST http://localhost:8080/relay \
  -H "Content-Type: application/json" \
  -d '{
    "ephemeralPublicKey": "0xaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "burnAddress": "0x742d35Cc6634C0532925a3b844Bc9e7595f42eD5"
  }'
```

### Expected Success Response
```json
{
  "success": true,
  "message": "Transaction submitted successfully",
  "txHash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"
}
```

## Environment Configuration

Update your `.env` file with:

```env
# Required configuration
PORT=8080
ETHEREUM_RPC=https://sepolia.infura.io/v3/YOUR_INFURA_KEY
CONTRACT_ADDRESS=0xYourContractAddressHere
PRIVATE_KEY=YourPrivateKeyWithoutPrefix

# Complete ABI from StealthBurnRegistry contract
CONTRACT_ABI=[{"inputs":[],"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"bytes32","name":"ephemeralPublicKey","type":"bytes32"},{"indexed":false,"internalType":"address","name":"burnAddress","type":"address"}],"name":"StealthBurn","type":"event"},{"inputs":[{"internalType":"bytes32","name":"ephemeralPublicKey","type":"bytes32"},{"internalType":"address","name":"burnAddress","type":"address"}],"name":"emitBurn","outputs":[],"stateMutability":"nonpayable","type":"function"}]
```

## Files Modified

1. **internal/ethereum/client.go**
   - Updated CallContract() function
   - Added parseBytes32() helper function
   - Added min() helper function
   - Added encoding/hex import

2. **.env.example**
   - Updated CONTRACT_ABI with complete StealthBurnRegistry ABI

3. **QUICK_START.md**
   - Updated API example with correct parameters
   - Added parameter format documentation

## Backward Compatibility

If you need to use a different contract:

1. Update `CONTRACT_ADDRESS` in `.env`
2. Update `CONTRACT_ABI` in `.env` with your contract's ABI
3. Modify the function name in `internal/ethereum/client.go` line ~90:
   ```go
   data, err := c.contractABI.Pack("yourFunctionName", param1, param2)
   ```
4. Adjust parameter types and parsing as needed

## Debugging

Enable detailed logging by checking console output:
```
Packed data for emitBurn: ephemeralPublicKey=..., burnAddress=...
Transaction sent with hash: 0x...
```

If you encounter errors:
- Check the CONTRACT_ADDRESS is correct
- Verify CONTRACT_ABI is valid JSON
- Ensure parameters match expected types
- Verify account has sufficient balance for gas
