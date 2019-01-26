package app

import (
	"encoding/base64"
	"encoding/json"
	"errors"
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
	}
}

func verifySessionToken(req *http.Request, oauthTokenStorage OAuthTokenStorage) (*sessionToken, error) {
	cookie, err := req.Cookie(sessionTokenCookieName)

	if err != nil {
		return nil, fmt.Errorf("failed to read session token cookie: %s", err)
	}

	data, err := base64.URLEncoding.DecodeString(cookie.Value)

	if err != nil {
		return nil, fmt.Errorf("failed to decode session token cookie: %s", err)
	}

	var stok sessionToken

	if err = json.Unmarshal(data, &stok); err != nil {
		return nil, fmt.Errorf("failed to decode session token cookie: %s", err)
	}

	oauthToken, err := oauthTokenStorage.GetOAuthToken(req.Context(), stok.Shop)

	if err != nil {
		return nil, fmt.Errorf("failed to check session token cookie: %s", err)
	}

	if oauthToken == nil {
		return nil, fmt.Errorf("unknown shop `%s`", stok.Shop)
	}

	if !oauthToken.Equal(stok.OAuthToken) {
		return nil, errors.New("invalid session token")
	}

	return &stok, nil
}
