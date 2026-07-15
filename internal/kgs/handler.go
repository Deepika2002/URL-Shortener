package kgs

import (
	"encoding/json"
	"net/http"
)

type handler struct {
	sfGen *SnowflakeGenerator
}

type kgsResponse struct {
	ID string `json:"id"`
}

type kgsError struct {
	Error string `json:"error"`
}

// NewHandler creates a new HTTP multiplexer and wires up the KGS endpoints.
func NewHandler(sfGen *SnowflakeGenerator) http.Handler {
	h := &handler{
		sfGen: sfGen,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", h.handleNextID)

	return mux
}

func (h *handler) handleNextID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(kgsError{Error: "Method Not Allowed"})
		return
	}

	id, err := h.sfGen.NextID()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(kgsError{Error: "Internal Server Error generating ID"})
		return
	}

	encodedID := EncodeBase62(id)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(kgsResponse{ID: encodedID})
}
