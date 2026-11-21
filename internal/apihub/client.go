package apihub

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Model represents a model from the APIHub storefront
type Model struct {
	Key           string `json:"key"`
	Name          string `json:"name"`
	Maker         string `json:"maker"`
	Type          string `json:"type"`
	ContextLength *int   `json:"context_length"`
	Offers        []Offer `json:"offers"`
}

type Offer struct {
	OfferID     int    `json:"offer_id"`
	Seller      Seller `json:"seller"`
	Model       string `json:"model"`
	Pricing     string `json:"pricing"`
	Input       string `json:"input"`
	Output      string `json:"output"`
	Cache       string `json:"cache"`
	CacheWrite  string `json:"cache_write"`
	Request     string `json:"request"`
	Variants    []interface{} `json:"variants"`
	Unavailable *string `json:"unavailable"`
	Country     string `json:"country"`
	Direct      Direct `json:"direct"`
	Trust       Trust  `json:"trust"`
}

type Seller struct {
	Slug       string `json:"slug"`
	Name       string `json:"name"`
	LegalName  string `json:"legal_name"`
}

type Direct struct {
	Method     string  `json:"method"`
	Path       string  `json:"path"`
	StatusPath *string `json:"status_path"`
	TaskID     *string `json:"task_id"`
}

type Trust struct {
	Labels      interface{} `json:"labels"`
	Quality     *float64    `json:"quality"`
	Verdict     string      `json:"verdict"`
}

type StorefrontResponse struct {
	Data []Model `json:"data"`
}

type Client struct {
	BaseURL string
	Token   string
	HTTP    *http.Client
}

func NewClient(baseURL, token string) *Client {
	return &Client{
		BaseURL: baseURL,
		Token:   token,
		HTTP: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// FetchModels fetches all models from the storefront
func (c *Client) FetchModels() ([]Model, error) {
	url := c.BaseURL + "/api/v1/storefront/models"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}
	if c.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.Token)
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.HTTP.Do(req)
	if err != nil {
		return nil, fmt.Errorf("fetching models: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	var result StorefrontResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	return result.Data, nil
}

// ChatModels returns only chat-type models
func (c *Client) ChatModels() ([]Model, error) {
	models, err := c.FetchModels()
	if err != nil {
		return nil, err
	}
	var chat []Model
	for _, m := range models {
		if m.Type == "chat" {
			chat = append(chat, m)
		}
	}
	return chat, nil
}
