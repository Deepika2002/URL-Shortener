package kgsclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"urlshortener/pkg/config"
)

// Client defines the interface for communicating with the Key Generation Service.
type Client interface {
	GetNextID() (string, error)
}

type kgsClient struct {
	baseURL    string
	httpClient *http.Client
}

// kgsResponse defines the expected JSON payload from the KGS.
type kgsResponse struct {
	ID string `json:"id"`
}

// NewClient creates and returns a new KGS HTTP client.
func NewClient(cfg *config.Config) Client {
	// For local development and docker-compose, we use localhost.
	// In a real production orchestrator (e.g. Kubernetes), this would be the internal DNS name of the KGS service.
	baseURL := fmt.Sprintf("http://localhost:%s", cfg.KGSPort)
	
	return &kgsClient{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 2 * time.Second, // strict timeout for microservice communication
		},
	}
}

// GetNextID fetches the next unique Snowflake ID from the KGS.
func (c *kgsClient) GetNextID() (string, error) {
	resp, err := c.httpClient.Get(c.baseURL + "/")
	if err != nil {
		return "", fmt.Errorf("failed to call KGS: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("KGS returned non-200 status: %d", resp.StatusCode)
	}

	var data kgsResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return "", fmt.Errorf("failed to decode KGS response: %w", err)
	}

	return data.ID, nil
}
