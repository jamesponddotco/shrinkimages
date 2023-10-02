// Package fetch provides a cache client that can fetch data from a URL.
package fetch

import (
	"context"
	"fmt"
	"io"
	"mime"
	"net/http"
	"net/url"
	"path"

	"git.sr.ht/~jamesponddotco/httpx-go"
	"git.sr.ht/~jamesponddotco/pagecache-go"
	"git.sr.ht/~jamesponddotco/pagecache-go/memorycachex"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/meta"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
	"golang.org/x/time/rate"
)

const (
	// ErrFetchData is returned when the client fails to fetch data from a URL.
	ErrFetchData xerrors.Error = "failed to fetch data"

	// ErrInvalidProtocol is returned when the protocol of a URL is invalid.
	ErrInvalidProtocol xerrors.Error = "invalid protocol; must be https"
)

// Client represents a client that can fetch data from a URL.
type Client struct {
	// httpc is the underlying HTTP client used to fetch data.
	httpc *httpx.Client
}

// New creates a new client that can fetch data from a URL.
func New(serviceName, serviceEmail string) *Client {
	return &Client{
		httpc: &httpx.Client{
			RateLimiter: rate.NewLimiter(rate.Limit(2), 1),
			RetryPolicy: httpx.DefaultRetryPolicy(),
			UserAgent: &httpx.UserAgent{
				Token:   serviceName,
				Version: meta.Version,
				Comment: []string{serviceEmail},
			},
			Cache: memorycachex.NewCache(pagecache.DefaultPolicy(), pagecache.DefaultCapacity),
		},
	}
}

// Remote fetches data from a URL, optimizes it to reduce its size, and returns
// it as a byte slice.
func (c *Client) Remote(ctx context.Context, uri string) (data []byte, filename string, err error) {
	u, err := url.Parse(uri)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrFetchData, err)
	}

	if u.Scheme != "https" {
		return nil, "", fmt.Errorf("%w: %s", ErrInvalidProtocol, u.Scheme)
	}

	resp, err := c.httpc.Get(ctx, uri)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrFetchData, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, "", fmt.Errorf("%w: %s", ErrFetchData, resp.Status)
	}

	data, err = io.ReadAll(resp.Body)
	if err != nil {
		return nil, "", fmt.Errorf("%w: %w", ErrFetchData, err)
	}

	filename = extractFilename(resp)

	return data, filename, nil
}

// extractFilename extracts the filename from the Content-Disposition header of
// an HTTP response. If not found, it falls back to the base name of the URL.
func extractFilename(resp *http.Response) string {
	_, params, err := mime.ParseMediaType(resp.Header.Get("Content-Disposition"))
	if err != nil || params["filename"] == "" {
		// Fallback to base name of the URL if there's an error or filename
		// isn't provided.
		uri, err := url.Parse(resp.Request.URL.String())
		if err != nil {
			return ""
		}

		return path.Base(uri.Path)
	}

	return params["filename"]
}
