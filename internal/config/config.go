// Package config implements the configuration logic for the Shrink Images
// service.
package config

import (
	"encoding/json"
	"fmt"
	"net/url"
	"os"

	"git.sr.ht/~jamesponddotco/shrinkimages/internal/meta"
	"git.sr.ht/~jamesponddotco/shrinkimages/internal/validate"
	"git.sr.ht/~jamesponddotco/xstd-go/xerrors"
)

const (
	// ErrInvalidConfigFile is returned when the configuration file is invalid.
	ErrInvalidConfigFile xerrors.Error = "invalid configuration file"

	// ErrMissingContact is returned when the contact information is missing.
	ErrMissingContact xerrors.Error = "service's contact information is missing"

	// ErrMissingPrivacyPolicy is returned when the privacy policy is missing.
	ErrMissingPrivacyPolicy xerrors.Error = "service's privacy policy is missing"

	// ErrMissingTermsOfService is returned when the terms of service is missing.
	ErrMissingTermsOfService xerrors.Error = "service's terms of service is missing"

	// ErrMissingAPIKey is returned when the API key is missing.
	ErrMissingAPIKey xerrors.Error = "service's API key is missing" //nolint:gosec // false positive

	// ErrMissingTLSCertificate is returned when the TLS certificate is missing.
	ErrMissingTLSCertificate xerrors.Error = "server's TLS certificate is missing"

	// ErrMissingTLSKey is returned when the TLS key is missing.
	ErrMissingTLSKey xerrors.Error = "server's TLS key is missing"

	// ErrInvalidTLSVersion is returned when the TLS version is invalid.
	ErrInvalidTLSVersion xerrors.Error = "server's TLS version is invalid; must be 1.2 or 1.3"

	// ErrInvalidHomepage is returned when the homepage is invalid.
	ErrInvalidHomepage xerrors.Error = "service's homepage is invalid"

	// ErrInvalidPrivacyPolicy is returned when the privacy policy is invalid.
	ErrInvalidPrivacyPolicy xerrors.Error = "service's privacy policy is invalid"

	// ErrInvalidTermsOfService is returned when the terms of service is invalid.
	ErrInvalidTermsOfService xerrors.Error = "service's terms of service is invalid"

	// ErrInvalidAPIKey is returned when the API key is invalid.
	ErrInvalidAPIKey xerrors.Error = "service's API key is invalid" //nolint:gosec // false positive
)

const (
	// DefaultMinTLSVersion is the default minimum TLS version supported by the
	// server.
	DefaultMinTLSVersion string = "1.3"

	// DefaultAddress is the default address of the application.
	DefaultAddress string = ":1997"

	// DefaultPID is the default path to the PID file.
	DefaultPID string = "/var/run/shrinkimages.pid"

	// DefaultServiceName is the default name of the service.
	DefaultServiceName string = meta.Name

	// DefaultHomepage is the default link to the service's homepage.
	DefaultHomepage string = meta.Homepage

	// DefaultMaxUploadSize is the default maximum upload size in megabytes.
	DefaultMaxUploadSize uint64 = 50

	// DefaultMaxAllowedWidth is the default maximum allowed width of the image.
	DefaultMaxAllowedWidth uint = 10000

	// DefaultMaxAllowedHeight is the default maximum allowed height of the
	// image.
	DefaultMaxAllowedHeight uint = 10000
)

// TLS represents the TLS configuration.
type TLS struct {
	// Certificate is the path to the TLS certificate.
	Certificate string `json:"certificate"`

	// Key is the path to the TLS key.
	Key string `json:"key"`

	// Version is the TLS version to use.
	Version string `json:"version"`
}

// Server represents the server configuration.
type Server struct {
	// TLS is the TLS configuration.
	TLS *TLS `json:"tls"`

	// Address is the address of the application.
	Address string `json:"address"`

	// PID is the path to the PID file.
	PID string `json:"pid"`
}

// Service represents the service configuration.
type Service struct {
	// Name is the name of the service.
	Name string `json:"name"`

	// Homepage is the link to the service's homepage.
	Homepage string `json:"homepage"`

	// Contact is the contact email for the service.
	Contact string `json:"contact"`

	// PrivacyPolicy is the link to the service's privacy policy.
	PrivacyPolicy string `json:"privacyPolicy"`

	// TermsOfService is the link to the service's terms of service.
	TermsOfService string `json:"termsOfService"`

	// APIKey is the API key for the service.
	APIKey string `json:"apiKey"`

	// MaxUploadSize is the maximum upload size in megabytes.
	MaxUploadSize uint64 `json:"maxUploadSize"`

	// MaxAllowedWidth is the maximum allowed width of the image.
	MaxAllowedWidth uint `json:"maxAllowedWidth"`

	// MaxAllowedHeight is the maximum allowed height of the image.
	MaxAllowedHeight uint `json:"maxAllowedHeight"`
}

// Config represents the application configuration.
type Config struct {
	// Service is the service configuration.
	Service *Service `json:"service"`

	// Server is the server configuration.
	Server *Server `json:"server"`
}

// LoadConfig opens a file and reads the configuration from it.
func LoadConfig(path string) (*Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}
	defer file.Close()

	config := &Config{
		Service: &Service{
			Name:             DefaultServiceName,
			Homepage:         DefaultHomepage,
			MaxUploadSize:    DefaultMaxUploadSize,
			MaxAllowedWidth:  DefaultMaxAllowedWidth,
			MaxAllowedHeight: DefaultMaxAllowedHeight,
		},
		Server: &Server{
			Address: DefaultAddress,
			PID:     DefaultPID,
			TLS: &TLS{
				Version: DefaultMinTLSVersion,
			},
		},
	}

	if err := json.NewDecoder(file).Decode(config); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("%w: %w", ErrInvalidConfigFile, err)
	}

	return config, nil
}

// Validate checks the configuration for errors.
func (c *Config) Validate() error { //nolint:gocyclo // TODO: find a better way to do this
	if c.Server == nil || c.Server.TLS == nil || c.Server.TLS.Certificate == "" {
		return ErrMissingTLSCertificate
	}

	if c.Server.TLS.Key == "" {
		return ErrMissingTLSKey
	}

	if c.Server.TLS.Version != "1.2" && c.Server.TLS.Version != "1.3" {
		return ErrInvalidTLSVersion
	}

	if c.Service == nil || c.Service.Contact == "" {
		return ErrMissingContact
	}

	if c.Service.PrivacyPolicy == "" {
		return ErrMissingPrivacyPolicy
	}

	if c.Service.TermsOfService == "" {
		return ErrMissingTermsOfService
	}

	if c.Service.APIKey == "" {
		return ErrMissingAPIKey
	}

	if _, err := url.Parse(c.Service.Homepage); err != nil {
		return ErrInvalidHomepage
	}

	if _, err := url.Parse(c.Service.PrivacyPolicy); err != nil {
		return ErrInvalidPrivacyPolicy
	}

	if _, err := url.Parse(c.Service.TermsOfService); err != nil {
		return ErrInvalidTermsOfService
	}

	if !validate.APIKey(c.Service.APIKey) {
		return ErrInvalidAPIKey
	}

	return nil
}
