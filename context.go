package shopify

import "context"

type contextKey int

const (
	contextKeyShop contextKey = iota
	contextKeyOAuthToken
)

// WithShop returns a context that references a shop.
func WithShop(ctx context.Context, shop Shop) context.Context {
	return context.WithValue(ctx, contextKeyShop, shop)
}

// GetShop returns the shop associated to a context.
func GetShop(ctx context.Context) (Shop, bool) {
	if v := ctx.Value(contextKeyShop); v != nil {
		return v.(Shop), true
	}

	return "", false
}

// WithOAuthToken returns a context that references an OAuth token.
//
// This will override any previously set access token with WithAccessToken.
func WithOAuthToken(ctx context.Context, token *OAuthToken) context.Context {
	return context.WithValue(ctx, contextKeyOAuthToken, token)
}

// GetOAuthToken returns the OAuth token associated to a context.
func GetOAuthToken(ctx context.Context) (*OAuthToken, bool) {
	if v := ctx.Value(contextKeyOAuthToken); v != nil {
		return v.(*OAuthToken), true
	}

	return nil, false
}

// WithAccessToken returns a context that references an access token.
//
// This will override any previously set OAuth token with WithOAuthToken.
func WithAccessToken(ctx context.Context, token AccessToken) context.Context {
	return WithOAuthToken(ctx, &OAuthToken{
		AccessToken: token,
	})
}

// GetAccessToken returns the access token associated to a context.
func GetAccessToken(ctx context.Context) (AccessToken, bool) {
	if oauthToken, ok := GetOAuthToken(ctx); ok {
		return oauthToken.AccessToken, true
	}

	return "", false
}
