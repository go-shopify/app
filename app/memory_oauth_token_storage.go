package app

import (
	"context"
	"sync"

	"github.com/go-shopify/shopify"
)

// MemoryOAuthTokenStorage implements in-memory storage of OAuth tokens.
type MemoryOAuthTokenStorage struct {
	oauthTokens map[shopify.Shop]shopify.OAuthToken
	lock        sync.Mutex
}

// GetOAuthToken gets an OAuth token for the specified shop.
//
// The method never fails.
//
// If no OAuth token exists for the shop, a nil OAuth token is returned.
func (s *MemoryOAuthTokenStorage) GetOAuthToken(ctx context.Context, shop shopify.Shop) (*shopify.OAuthToken, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if oauthToken, ok := s.dict()[shop]; ok {
		return &oauthToken, nil
	}

	return nil, nil
}

// UpdateOAuthToken updates an OAuth token.
//
// If the shop has no previous OAuth token, it is then created.
//
// The method never fails.
func (s *MemoryOAuthTokenStorage) UpdateOAuthToken(ctx context.Context, shop shopify.Shop, oauthToken shopify.OAuthToken) error {
	s.lock.Lock()

	s.dict()[shop] = oauthToken

	s.lock.Unlock()

	return nil
}

// DeleteOAuthToken deletes an OAuth token for a shop.
//
// If the shop has no OAuth token, the call is a no-op.
//
// The method never fails.
func (s *MemoryOAuthTokenStorage) DeleteOAuthToken(ctx context.Context, shop shopify.Shop) error {
	s.lock.Lock()

	delete(s.dict(), shop)

	s.lock.Unlock()

	return nil
}

func (s *MemoryOAuthTokenStorage) dict() map[shopify.Shop]shopify.OAuthToken {
	if s.oauthTokens == nil {
		s.oauthTokens = map[shopify.Shop]shopify.OAuthToken{}
	}

	return s.oauthTokens
}
