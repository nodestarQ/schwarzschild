package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port            string
	EthereumRPC     string
	ContractAddress string
	PrivateKey      string
	ContractABI     string
}

func Load() (*Config, error) {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Port:            getEnv("PORT", "8080"),
		EthereumRPC:     getEnv("ETHEREUM_RPC", "https://sepolia.infura.io/v3/YOUR_INFURA_KEY"),
		ContractAddress: getEnv("CONTRACT_ADDRESS", ""),
		PrivateKey:      getEnv("PRIVATE_KEY", ""),
		ContractABI:     getEnv("CONTRACT_ABI", ""),
	}

	// Validate required fields
	if cfg.ContractAddress == "" {
		return nil, fmt.Errorf("CONTRACT_ADDRESS is required")
	}

	if cfg.PrivateKey == "" {
		return nil, fmt.Errorf("PRIVATE_KEY is required")
	}

	if cfg.ContractABI == "" {
		return nil, fmt.Errorf("CONTRACT_ABI is required")
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
