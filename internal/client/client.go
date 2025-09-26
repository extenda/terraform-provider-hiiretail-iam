package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// IClient is an interface for the Client type
type IClient interface {
	CreateGroup(ctx context.Context, name string, description string) (*Group, error)
	GetGroup(ctx context.Context, id string) (*Group, error)
	UpdateGroup(ctx context.Context, id string, name string, description string) (*Group, error)
	DeleteGroup(ctx context.Context, id string) error
}

// Client represents the HiiRetail IAM API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new HiiRetail IAM API client
func NewClient(baseURL string, token string) *Client {
	return &Client{
		baseURL:    strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{},
		token:      token,
	}
}

// Group represents an IAM group
type Group struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// CreateGroup creates a new IAM group
func (c *Client) CreateGroup(ctx context.Context, name string, description string) (*Group, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group payload: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/groups", strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	var group Group
	if err := c.do(req, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// GetGroup retrieves an IAM group by ID
func (c *Client) GetGroup(ctx context.Context, id string) (*Group, error) {
	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/groups/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var group Group
	if err := c.do(req, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// UpdateGroup updates an existing IAM group
func (c *Client) UpdateGroup(ctx context.Context, id string, name string, description string) (*Group, error) {
	payload := map[string]interface{}{
		"name":        name,
		"description": description,
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal group payload: %w", err)
	}

	req, err := c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/groups/%s", id), strings.NewReader(string(data)))
	if err != nil {
		return nil, err
	}

	var group Group
	if err := c.do(req, &group); err != nil {
		return nil, err
	}

	return &group, nil
}

// DeleteGroup deletes an IAM group
func (c *Client) DeleteGroup(ctx context.Context, id string) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/groups/%s", id), nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}

// newRequest creates a new HTTP request with common headers
func (c *Client) newRequest(ctx context.Context, method string, path string, body io.Reader) (*http.Request, error) {
	req, err := http.NewRequestWithContext(ctx, method, c.baseURL+path, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	return req, nil
}

// do performs an HTTP request and decodes the response
func (c *Client) do(req *http.Request, v interface{}) error {
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
	}

	if v != nil {
		if err := json.NewDecoder(resp.Body).Decode(v); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}