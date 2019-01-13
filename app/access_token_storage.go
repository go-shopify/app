package app

import (
	"context"
	"sync"

	"github.com/go-shopify/shopify"
)

// AccessTokenStorage represents an access token storage.
type AccessTokenStorage interface {
	// GetAccessToken gets an access token for the specified shop.
	//
	// If the request fails, an error is returned.
	//
	// If no access token exists for the shop, an empty access token is
	// returned.
	GetAccessToken(ctx context.Context, shop shopify.Shop) (shopify.AccessToken, error)

	// UpdateAccessToken updates an access token.
	//
	// If the shop has no previous access token, it is then created.
	UpdateAccessToken(ctx context.Context, shop shopify.Shop, accessToken shopify.AccessToken) error

	// DeleteAccessToken deletes an access token for a shop.
	//
	// If the shop has no access token, the call is a no-op.
	DeleteAccessToken(ctx context.Context, shop shopify.Shop) error
}

// MemoryAccessTokenStorage implements in-memory storage of access tokens.
type MemoryAccessTokenStorage struct {
	accessTokens map[shopify.Shop]shopify.AccessToken
	lock         sync.Mutex
}

// GetAccessToken gets an access token for the specified shop.
//
// The method never fails.
//
// If no access token exists for the shop, an empty access token is
// returned.
func (s *MemoryAccessTokenStorage) GetAccessToken(ctx context.Context, shop shopify.Shop) (shopify.AccessToken, error) {
	s.lock.Lock()

	accessToken, _ := s.dict()[shop]

	s.lock.Unlock()

	return accessToken, nil
}

// UpdateAccessToken updates an access token.
//
// If the shop has no previous access token, it is then created.
//
// The method never fails.
func (s *MemoryAccessTokenStorage) UpdateAccessToken(ctx context.Context, shop shopify.Shop, accessToken shopify.AccessToken) error {
	s.lock.Lock()

	s.dict()[shop] = accessToken

	s.lock.Unlock()

	return nil
}

// DeleteAccessToken deletes an access token for a shop.
//
// If the shop has no access token, the call is a no-op.
//
// The method never fails.
func (s *MemoryAccessTokenStorage) DeleteAccessToken(ctx context.Context, shop shopify.Shop) error {
	s.lock.Lock()

	delete(s.dict(), shop)

	s.lock.Unlock()

	return nil
}

func (s *MemoryAccessTokenStorage) dict() map[shopify.Shop]shopify.AccessToken {
	if s.accessTokens == nil {
		s.accessTokens = map[shopify.Shop]shopify.AccessToken{}
	}

	return s.accessTokens
}
