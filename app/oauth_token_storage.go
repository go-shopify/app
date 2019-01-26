package app

import (
	"context"

	"github.com/go-shopify/shopify"
)

// OAuthTokenStorage represents an OAuth token storage.
type OAuthTokenStorage interface {
	// GetOAuthToken gets an OAuth token for the specified shop.
	//
	// If the request fails, an error is returned.
	//
	// If no OAuth token exists for the shop, a nil OAuth token is returned.
	GetOAuthToken(ctx context.Context, shop shopify.Shop) (*shopify.OAuthToken, error)

	// UpdateOAuthToken updates an OAuth token.
	//
	// If the shop has no previous OAuth token, it is then created.
	UpdateOAuthToken(ctx context.Context, shop shopify.Shop, oauthToken shopify.OAuthToken) error

	// DeleteOAuthToken deletes an OAuth token for a shop.
	//
	// If the shop has no OAuth token, the call is a no-op.
	DeleteOAuthToken(ctx context.Context, shop shopify.Shop) error
}
