// Package cryptoutil provides utility functions for cryptographic operations.
package cryptoutil

import (
	"encoding/base64"
	"fmt"
	"io"
)

// KeyPrefix is the prefix for all keys generated by this package.
const KeyPrefix string = "shrink-"

// GenerateAPIKey generates a new API key with the given random source.
func GenerateAPIKey(random io.Reader) (string, error) {
	key := make([]byte, 24) //nolint:makezero // we don't need a zeroed initial length

	_, err := io.ReadFull(random, key)
	if err != nil {
		return "", fmt.Errorf("failed to read random source: %w", err)
	}

	return KeyPrefix + base64.StdEncoding.EncodeToString(key), nil
}
