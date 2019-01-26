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

// AuthenticatedHandler represents a HTTP handler that contains additional shop and OAuth information.
type AuthenticatedHandler interface {
	ServeHTTPAuthenticated(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken)
}

// The AuthenticatedHandlerFunc type is an adapter to allow the use of ordinary functions as authenticated HTTP handlers.
//
// If f is a function with the appropriate signature, AuthenticatedHandlerFunc(f) is an AuthenticatedHandler that calls f.
type AuthenticatedHandlerFunc func(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken)

// ServeHTTPAuthenticated calls f(w, req, shop, oauthToken).
func (f AuthenticatedHandlerFunc) ServeHTTPAuthenticated(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken) {
	f(w, req, shop, oauthToken)
}

// NewSessionHandler instantiates a new session handler.
func NewSessionHandler(authenticatedHandler AuthenticatedHandler, oauthTokenStorage OAuthTokenStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		stok, err := verifySessionToken(req, oauthTokenStorage)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Error: %s.", err)
			return
		}

		// Make sure to refresh the cookie.
		http.SetCookie(w, stok.AsCookie())

		authenticatedHandler.ServeHTTPAuthenticated(w, req, stok.Shop, &stok.OAuthToken)
	})
}
