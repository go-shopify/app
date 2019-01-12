package app

import "net/http"

// Config contains the configuration for a Handler.
type Config struct {
	// APIKey is the Shopify API key for the app, as indicated on the Shopify
	// App page.
	APIKey string

	// APISecret is the Shopify API secret for the app, as indicated on the
	// Shopify App page.
	APISecret string

	// PublicURL is the public URL at which the app will be instantiated.
	PublicURL string

	// The Scope of the app, as documented at
	// https://help.shopify.com/en/api/getting-started/authentication/oauth/scopes.
	Scope []string

	// DefaultHandler is the handler to call when no known route was matched.
	//
	// If none is specified, a 404 is returned for unknown routes.
	DefaultHandler http.Handler
}
