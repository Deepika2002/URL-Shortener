package read

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

type handler struct {
	svc Service
}

type errorResponse struct {
	Error string `json:"error"`
}

// NewHandler creates a new HTTP multiplexer and wires up the Read Service endpoints.
func NewHandler(svc Service) http.Handler {
	h := &handler{
		svc: svc,
	}

	mux := http.NewServeMux()
	// Catch-all route to handle dynamic /{shortID} paths
	mux.HandleFunc("/", h.handleRedirect)
	return mux
}

func (h *handler) handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, "Method Not Allowed")
		return
	}

	// Extract the shortID from the path, stripping the leading "/"
	shortID := strings.TrimPrefix(r.URL.Path, "/")
	if shortID == "" {
		writeError(w, http.StatusBadRequest, "Short ID is required")
		return
	}

	// Attempt to retrieve the original URL
	longURL, err := h.svc.GetLongURL(r.Context(), shortID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			writeError(w, http.StatusNotFound, "Short URL not found")
			return
		}
		writeError(w, http.StatusInternalServerError, "Internal server error while resolving URL")
		return
	}

	// Perform the HTTP Redirect.
	// We use 302 Found (Temporary Redirect) instead of 301 (Permanent) to prevent browsers 
	// from heavily caching the result. This ensures every click hits our service so we can 
	// track accurate analytics.
	http.Redirect(w, r, longURL, http.StatusFound)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(errorResponse{Error: msg})
}
