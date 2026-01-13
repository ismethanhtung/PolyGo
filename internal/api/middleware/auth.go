package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/polygo/internal/config"
	"github.com/polygo/pkg/response"
)

// AuthCredentials holds extracted auth credentials from request
type AuthCredentials struct {
	APIKey     string
	APISecret  string
	Passphrase string
	Signature  string
	Timestamp  string
}

// AuthConfig holds auth middleware configuration
type AuthConfig struct {
	Config *config.AuthConfig
	// Optional: validate credentials against a store
	Validator func(creds *AuthCredentials) bool
}

// Auth returns a middleware that extracts and validates auth credentials
func Auth(cfg *config.AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		creds := &AuthCredentials{
			APIKey:     c.Get(cfg.APIKeyHeader),
			APISecret:  c.Get(cfg.APISecretHeader),
			Passphrase: c.Get(cfg.PassphraseHeader),
			Signature:  c.Get(cfg.SignatureHeader),
			Timestamp:  c.Get(cfg.TimestampHeader),
		}
		
		// Check required fields for authenticated endpoints
		if creds.APIKey == "" {
			return response.Unauthorized(c, "API key is required")
		}
		
		if creds.Timestamp == "" {
			return response.Unauthorized(c, "Timestamp is required")
		}
		
		if creds.Signature == "" {
			return response.Unauthorized(c, "Signature is required")
		}
		
		// Store credentials in context for handlers
		c.Locals("auth", creds)
		
		return c.Next()
	}
}

// OptionalAuth extracts auth credentials if present, but doesn't require them
func OptionalAuth(cfg *config.AuthConfig) fiber.Handler {
	return func(c *fiber.Ctx) error {
		apiKey := c.Get(cfg.APIKeyHeader)
		
		if apiKey != "" {
			creds := &AuthCredentials{
				APIKey:     apiKey,
				APISecret:  c.Get(cfg.APISecretHeader),
				Passphrase: c.Get(cfg.PassphraseHeader),
				Signature:  c.Get(cfg.SignatureHeader),
				Timestamp:  c.Get(cfg.TimestampHeader),
			}
			c.Locals("auth", creds)
		}
		
		return c.Next()
	}
}

// GetAuthCredentials retrieves auth credentials from context
func GetAuthCredentials(c *fiber.Ctx) *AuthCredentials {
	if creds, ok := c.Locals("auth").(*AuthCredentials); ok {
		return creds
	}
	return nil
}

// GetAuthHeaders converts credentials to headers map for upstream requests
func GetAuthHeaders(creds *AuthCredentials, cfg *config.AuthConfig) map[string]string {
	if creds == nil {
		return nil
	}
	
	return map[string]string{
		cfg.APIKeyHeader:     creds.APIKey,
		cfg.APISecretHeader:  creds.APISecret,
		cfg.PassphraseHeader: creds.Passphrase,
		cfg.SignatureHeader:  creds.Signature,
		cfg.TimestampHeader:  creds.Timestamp,
	}
}

// RequireAuth is a simple middleware that checks if auth credentials are present
func RequireAuth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if GetAuthCredentials(c) == nil {
			return response.Unauthorized(c, "Authentication required")
		}
		return c.Next()
	}
}
