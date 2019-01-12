package app

import (
	"context"
	"net/http"
	"net/url"

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

	// The Scopes of the app, as documented at
	// https://help.shopify.com/en/api/getting-started/authentication/oauth/scopes.
	Scopes shopify.Scopes

	// Handler is the handler to defer normal requests to.
	Handler http.Handler

	// OnAccessTokenRequested is a function to call whenever an access token is
	// requested for a given shop.
	//
	// If an error is returned, the request fails.

	// If an empty access token is returned, the OAuth authentication cycle
	// will start.
	OnAccessTokenRequested func(ctx context.Context, shop shopify.Shop) (shopify.AccessToken, error)

	// OnAccessTokenUpdated is a function to call whenever an access token for
	// a shop was updated.
	OnAccessTokenUpdated func(ctx context.Context, shop shopify.Shop, accessToken shopify.AccessToken) error

	// OnAccessTokenDeleted is a function to call whenever an access token for
	// a shop should be deleted.
	OnAccessTokenDelete func(ctx context.Context, shop shopify.Shop) error
}
