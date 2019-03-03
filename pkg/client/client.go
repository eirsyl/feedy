package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/eirsyl/feedy/pkg"
)

// Client for communicating with the gateway
type Client struct {
	client    *http.Client
	UserAgent string
}

// NewHTTPClient creates a new http client.
func NewHTTPClient() (*http.Client, error) {
	client := &http.Client{Timeout: ClientTimeout}
	return client, nil
}

// New creates a new client
func New() (*Client, error) {
	httpClient, err := NewHTTPClient()
	if err != nil {
		return nil, err
	}

	ua := fmt.Sprintf("Feedy %s", pkg.Version)
	c := &Client{
		client:    httpClient,
		UserAgent: ua,
	}

	return c, nil
}

// NewRequest creates a new http request
func (c *Client) NewRequest(method, url string, body interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if body != nil {
		buf = new(bytes.Buffer)
		enc := json.NewEncoder(buf)
		enc.SetEscapeHTML(false)
		err := enc.Encode(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, url, buf)
	if err != nil {
		return nil, err
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json; charset=UTF8")
	}
	//req.Header.Set("Accept", "application/json")
	req.Header.Set("X-Accept", "application/json")

	if c.UserAgent != "" {
		req.Header.Set("User-Agent", c.UserAgent)
	}
	return req, nil
}

// Do executes an http request and returns a response
// nolint: gocyclo
func (c *Client) Do(ctx context.Context, req *http.Request, v interface{}) (*http.Response, error) {
	req = withContext(ctx, req)

	resp, err := c.client.Do(req)
	if err != nil {
		// If we got an error, and the context has been canceled,
		// the context's error is probably more useful.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		// If the error type is *url.Error, sanitize its URL before returning.
		if e, ok := err.(*url.Error); ok {
			var url *url.URL
			if url, err = url.Parse(e.URL); err == nil {
				e.URL = sanitizeURL(url).String()
				return nil, e
			}
		}

		return nil, err
	}
	defer resp.Body.Close() // nolint: errcheck

	err = CheckResponse(resp)
	if err != nil {
		return resp, err
	}

	if v != nil {
		if w, ok := v.(io.Writer); ok {
			_, err = io.Copy(w, resp.Body)
			if err != nil {
				return nil, err
			}
		} else {
			decErr := json.NewDecoder(resp.Body).Decode(v)
			if decErr == io.EOF {
				decErr = nil // ignore EOF errors caused by empty response body
			}
			if decErr != nil {
				err = decErr
			}
		}
	}

	return resp, err
}

func withContext(ctx context.Context, req *http.Request) *http.Request {
	return req.WithContext(ctx)
}

func sanitizeURL(uri *url.URL) *url.URL {
	if uri == nil {
		return nil
	}

	return uri
}

// CheckResponse validates the response and raises an error if necceceary
func CheckResponse(r *http.Response) error {
	// Check status code
	if !(r.StatusCode >= 200 && r.StatusCode < 300) {
		return fmt.Errorf("Response status not in range [200, 300], actual code %d", r.StatusCode)
	}

	return nil
}
