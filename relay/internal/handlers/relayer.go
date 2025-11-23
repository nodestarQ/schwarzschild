package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/schwarzschild/relay/internal/config"
	"github.com/schwarzschild/relay/internal/ethereum"
)

type RelayRequest struct {
	EphemeralPublicKey string `json:"ephemeralPublicKey"`
	BurnAddress        string `json:"burnAddress"`
}

type RelayResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	TxHash  string `json:"txHash,omitempty"`
	Error   string `json:"error,omitempty"`
}

type RelayHandler struct {
	ethClient *ethereum.Client
}

func NewRelayHandler(cfg *config.Config) *RelayHandler {
	ethClient := ethereum.NewClient(cfg)
	return &RelayHandler{
		ethClient: ethClient,
	}
}

func (h *RelayHandler) HandleRelay(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight requests
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	var req RelayRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RelayResponse{
			Success: false,
			Error:   "Invalid request body: " + err.Error(),
		})
		return
	}

	// Validate input
	if req.EphemeralPublicKey == "" || req.BurnAddress == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RelayResponse{
			Success: false,
			Error:   "ephemeralPublicKey and burnAddress are required",
		})
		return
	}

	log.Printf("Relay request received - PublicKey: %s, BurnAddress: %s", req.EphemeralPublicKey, req.BurnAddress)

	// Call contract with the request data
	txHash, err := h.ethClient.CallContract(req.EphemeralPublicKey, req.BurnAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RelayResponse{
			Success: false,
			Error:   "Failed to relay transaction: " + err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RelayResponse{
		Success: true,
		Message: "Transaction submitted successfully",
		TxHash:  txHash,
	})
}
