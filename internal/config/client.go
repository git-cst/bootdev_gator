package config

import (
	"net/http"
	"time"
)

type ClientOptions struct {
	Timeout     time.Duration
	UserAgent   string
	ContentType string
	Headers     map[string]string
}

type Client struct {
	Http        http.Client
	UserAgent   string
	ContentType string
	Headers     map[string]string
}

func NewClient(options ClientOptions) Client {
	return Client{
		Http: http.Client{
			Timeout: options.Timeout,
		},
		UserAgent:   options.UserAgent,
		ContentType: options.ContentType,
		Headers:     options.Headers,
	}
}

// Do performs an HTTP request and applies all configured headers
func (c *Client) Do(req *http.Request) (*http.Response, error) {
	// Set the User-Agent if configured
	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}

	// Set the Content-Type if configured
	if c.ContentType != "" {
		req.Header.Set("Content-Type", c.ContentType)
	}

	// Add any additional custom headers
	for key, value := range c.Headers {
		req.Header.Set(key, value)
	}

	// Use the underlying http.Client to perform the request
	return c.Http.Do(req)
}
