package ethereum

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/schwarzschild/relay/internal/config"
)

type Client struct {
	rpcClient       *ethclient.Client
	privateKey      *ecdsa.PrivateKey
	contractAddress common.Address
	contractABI     *abi.ABI
	config          *config.Config
}

func NewClient(cfg *config.Config) *Client {
	// Connect to Ethereum RPC
	client, err := ethclient.Dial(cfg.EthereumRPC)
	if err != nil {
		log.Fatalf("Failed to connect to Ethereum: %v", err)
	}

	// Parse private key
	privateKey, err := crypto.HexToECDSA(cfg.PrivateKey)
	if err != nil {
		log.Fatalf("Failed to parse private key: %v", err)
	}

	// Parse contract ABI
	contractABI, err := abi.JSON(strings.NewReader(cfg.ContractABI))
	if err != nil {
		log.Fatalf("Failed to parse contract ABI: %v", err)
	}

	contractAddr := common.HexToAddress(cfg.ContractAddress)

	return &Client{
		rpcClient:       client,
		privateKey:      privateKey,
		contractAddress: contractAddr,
		contractABI:     &contractABI,
		config:          cfg,
	}
}

func (c *Client) CallContract(ephemeralPublicKey, burnAddress string) (string, error) {
	ctx := context.Background()

	// Get the public address from private key
	fromAddress := crypto.PubkeyToAddress(c.privateKey.PublicKey)

	// Get nonce for the transaction
	nonce, err := c.rpcClient.PendingNonceAt(ctx, fromAddress)
	if err != nil {
		return "", fmt.Errorf("failed to get nonce: %v", err)
	}

	// Get gas price
	gasPrice, err := c.rpcClient.SuggestGasPrice(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get gas price: %v", err)
	}

	// Create transaction opts
	auth, err := bind.NewKeyedTransactorWithChainID(c.privateKey, big.NewInt(11155111)) // Sepolia chain ID
	if err != nil {
		return "", fmt.Errorf("failed to create transactor: %v", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.Value = big.NewInt(0)
	auth.GasLimit = uint64(100000) // Adjusted for emitBurn function
	auth.GasPrice = gasPrice

	// Parse ephemeralPublicKey as bytes32
	// If it's a hex string, strip 0x prefix and pad to 32 bytes
	ephemeralBytes := parseBytes32(ephemeralPublicKey)

	// Parse burn address
	burnAddr := common.HexToAddress(burnAddress)

	// Pack the function call for emitBurn(bytes32 ephemeralPublicKey, address burnAddress)
	data, err := c.contractABI.Pack("emitBurn", ephemeralBytes, burnAddr)
	if err != nil {
		return "", fmt.Errorf("failed to pack contract call: %v", err)
	}

	log.Printf("Packed data for emitBurn: ephemeralPublicKey=%x, burnAddress=%s", ephemeralBytes, burnAddr.Hex())

	// Create the transaction
	tx := types.NewTx(&types.LegacyTx{
		Nonce:    nonce,
		GasPrice: gasPrice,
		Gas:      auth.GasLimit,
		To:       &c.contractAddress,
		Data:     data,
	})

	// Sign transaction
	signedTx, err := auth.Signer(auth.From, tx)
	if err != nil {
		return "", fmt.Errorf("failed to sign transaction: %v", err)
	}

	// Send transaction
	err = c.rpcClient.SendTransaction(ctx, signedTx)
	if err != nil {
		return "", fmt.Errorf("failed to send transaction: %v", err)
	}

	log.Printf("Transaction sent with hash: %s", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

// parseBytes32 converts a hex string to a [32]byte array
// Handles both with and without 0x prefix, and pads shorter values with zeros
func parseBytes32(hexStr string) [32]byte {
	var result [32]byte

	// Remove 0x prefix if present
	if len(hexStr) >= 2 && hexStr[:2] == "0x" {
		hexStr = hexStr[2:]
	}

	// If the hex string is shorter than 64 characters, pad with leading zeros
	if len(hexStr) < 64 {
		hexStr = strings.Repeat("0", 64-len(hexStr)) + hexStr
	}

	// Decode hex string to bytes
	decoded, err := hex.DecodeString(hexStr)
	if err != nil {
		log.Printf("Warning: failed to parse bytes32, using zero value: %v", err)
		return result
	}

	// Copy decoded bytes into result array (only first 32 bytes)
	copy(result[:], decoded[:min(len(decoded), 32)])

	return result
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
