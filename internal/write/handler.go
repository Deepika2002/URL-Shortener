package write

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

type handler struct {
	svc     Service
	baseURL string
}

type shortenRequest struct {
	LongURL     string `json:"long_url"`
	CustomAlias string `json:"custom_alias,omitempty"`
}

type shortenResponse struct {
	ShortURL string `json:"short_url"`
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewHandler creates a new HTTP multiplexer and wires up the Write Service endpoints.
func NewHandler(svc Service, baseURL string) http.Handler {
	h := &handler{
		svc:     svc,
		baseURL: strings.TrimRight(baseURL, "/"),
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/shorten", h.handleShorten)
	return mux
}

func (h *handler) handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	var req shortenRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "Invalid JSON payload")
		return
	}

	// Basic URL validation
	if req.LongURL == "" {
		writeError(w, http.StatusBadRequest, "long_url is required")
		return
	}
	if _, err := url.ParseRequestURI(req.LongURL); err != nil {
		writeError(w, http.StatusBadRequest, "long_url must be a valid URL")
		return
	}

	// Delegate to the core service
	shortID, err := h.svc.ShortenURL(r.Context(), req.LongURL, req.CustomAlias)
	if err != nil {
		status := http.StatusInternalServerError
		// Check if it's an alias collision
		if strings.Contains(err.Error(), "already in use") {
			status = http.StatusConflict
		}
		writeError(w, status, err.Error())
		return
	}

	// Construct the final shortened URL
	finalURL := fmt.Sprintf("%s/%s", h.baseURL, shortID)
	resp := shortenResponse{
		ShortURL: finalURL,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
