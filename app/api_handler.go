package app

import (
	"fmt"
	"net/http"

	"github.com/go-shopify/shopify"
)

// AuthenticatedAPIHandler represents a HTTP handler that contains additional shop and OAuth information.
type AuthenticatedAPIHandler interface {
	ServeHTTPAPI(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken)
}

// The AuthenticatedAPIHandlerFunc type is an adapter to allow the use of ordinary functions as authenticated HTTP handlers.
//
// If f is a function with the appropriate signature, AuthenticatedAPIHandlerFunc(f) is an AuthenticatedAPIHandler that calls f.
type AuthenticatedAPIHandlerFunc func(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken)

// ServeHTTPAPI calls f(w, req, shop, oauthToken).
func (f AuthenticatedAPIHandlerFunc) ServeHTTPAPI(w http.ResponseWriter, req *http.Request, shop shopify.Shop, oauthToken *shopify.OAuthToken) {
	f(w, req, shop, oauthToken)
}

// NewAPIHandler instantiates a new API handler.
//
// A typical usage is to wrap custom API rest endpoints with an APIHandler to
// ensure that the calls originates from a Shopify admin page that went through
// a OAuthHandler.
func NewAPIHandler(handler AuthenticatedAPIHandler, oauthTokenStorage OAuthTokenStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		stok, err := verifySessionToken(req, oauthTokenStorage)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Error: %s.", err)
			return
		}

		// Make sure to refresh the cookie.
		http.SetCookie(w, stok.AsCookie())

		handler.ServeHTTPAPI(w, req, stok.Shop, &stok.OAuthToken)
	})
}
