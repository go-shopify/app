package app

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/go-shopify/shopify"
)

// Config contains the configuration for a Handler.
type Config struct {
	// APIKey is the Shopify API key for the app, as indicated on the Shopify
	// App page.
	APIKey shopify.APIKey

	// APISecret is the Shopify API secret for the app, as indicated on the
	// Shopify App page.
	APISecret shopify.APISecret

	// PublicURL is the public URL at which the app will be instantiated.
	PublicURL *url.URL

	// The Scope of the app, as documented at
	// https://help.shopify.com/en/api/getting-started/authentication/oauth/scopes.
	Scope shopify.Scope

	// OnError is a function to call whenever an error happens.
	OnError func(ctx context.Context, err error)

	// AuthCallbackPath is the path for the authentication callback.
	//
	// If not explicity set, a sane default is chosen.
	AuthCallbackPath string
}

const defaultAuthCallbackPath = "/auth/callback"

func (c *Config) authCallbackPath() string {
	if c.AuthCallbackPath == "" {
		return defaultAuthCallbackPath
	}

	return c.AuthCallbackPath
}

const (
	envShopifyAPIKey    = "SHOPIFY_API_KEY"
	envShopifyAPISecret = "SHOPIFY_API_SECRET"
	envShopifyPublicURL = "SHOPIFY_PUBLIC_URL"
	envShopifyScope     = "SHOPIFY_SCOPE"
)

// ReadConfigFromEnvironment reads a configuration from environment variables.
func ReadConfigFromEnvironment() (*Config, error) {
	publicURL, err := url.Parse(os.Getenv(envShopifyPublicURL))

	if err != nil {
		return nil, fmt.Errorf("incorrect `%s`: %s", envShopifyPublicURL, err)
	}

	scope, err := shopify.ParseScope(os.Getenv(envShopifyScope))

	if err != nil {
		return nil, fmt.Errorf("incorrect `%s`: %s", envShopifyScope, err)
	}

	config := &Config{
		APIKey:    shopify.APIKey(os.Getenv(envShopifyAPIKey)),
		APISecret: shopify.APISecret(os.Getenv(envShopifyAPISecret)),
		PublicURL: publicURL,
		Scope:     scope,
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("`%s` is not set", envShopifyAPIKey)
	}

	if config.APISecret == "" {
		return nil, fmt.Errorf("`%s` is not set", envShopifyAPISecret)
	}

	return config, nil
}
