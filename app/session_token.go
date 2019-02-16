package app

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-shopify/shopify"
)

const sessionTokenCookieName = "session-token"

type sessionToken struct {
	Shop       shopify.Shop       `json:"shop"`
	OAuthToken shopify.OAuthToken `json:"oauth_token"`
}

func (t sessionToken) AsCookie() *http.Cookie {
	data, _ := json.Marshal(&t)

	return &http.Cookie{
		Name:   sessionTokenCookieName,
		Value:  base64.URLEncoding.EncodeToString(data),
		Secure: true,
		MaxAge: 3600, // 1 hour.
	}
}

func (t *sessionToken) FromCookie(cookie *http.Cookie) error {
	data, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return fmt.Errorf("failed to decode session token cookie: %s", err)
	}

	if err = json.Unmarshal(data, t); err != nil {
		return fmt.Errorf("failed to decode session token cookie: %s", err)
	}

	return nil
}

func verifySessionToken(req *http.Request, oauthTokenStorage OAuthTokenStorage) (*sessionToken, error) {
	shopCookie, err := req.Cookie(shopifyShopCookieName)

	if err != nil {
		return nil, fmt.Errorf("failed to read shop cookie: %s", err)
	}

	shop := shopify.Shop(shopCookie.Value)

	oauthToken, err := oauthTokenStorage.GetOAuthToken(req.Context(), shop)

	if err != nil {
		return nil, fmt.Errorf("failed to check session token cookie: %s", err)
	}

	if oauthToken == nil {
		return nil, fmt.Errorf("unknown shop `%s`", shop)
	}

	for _, cookie := range req.Cookies() {
		if cookie.Name != sessionTokenCookieName {
			continue
		}

		var stok sessionToken

		if err := stok.FromCookie(cookie); err != nil {
			// Malformed cookie. Skip.
			continue
		}

		if stok.Shop != shop {
			// The cookie is for a different shop. Skip.
			continue
		}

		if !oauthToken.Equal(stok.OAuthToken) {
			continue
		}

		return &stok, nil
	}

	return nil, fmt.Errorf("missing session token cookie: %s", sessionTokenCookieName)
}

func withSessionToken(ctx context.Context, stok *sessionToken) context.Context {
	ctx = shopify.WithShop(ctx, stok.Shop)
	ctx = shopify.WithOAuthToken(ctx, &stok.OAuthToken)

	return ctx
}
