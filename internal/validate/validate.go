// Package validate provides a set of validation functions for use with the application.
package validate

import (
	"encoding/base64"
	"strings"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/cryptoutil"
)

const (
	// KeyPrefix is the prefix used for valid API keys.
	KeyPrefix string = cryptoutil.KeyPrefix

	// MinKeyLength is the minimum length of a valid API key.
	MinKeyLength int = 32
)

// APIKey checks a string to ensure it's a valid API key.
//
// A valid API key is a string prefixed with "shrink-" followed by a base64
// encoded string of at least 32 characters.
func APIKey(key string) bool {
	if !strings.HasPrefix(key, KeyPrefix) {
		return false
	}

	base64Key := strings.TrimPrefix(key, KeyPrefix)

	if len(base64Key) < MinKeyLength {
		return false
	}

	_, err := base64.RawStdEncoding.DecodeString(base64Key)

	return err == nil
}
