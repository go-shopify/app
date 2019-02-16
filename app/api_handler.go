package app

import (
	"fmt"
	"net/http"
)

// NewAPIHandler instantiates a new API handler.
//
// A typical usage is to wrap custom API rest endpoints with an APIHandler to
// ensure that the calls originates from a Shopify admin page that went through
// a OAuthHandler.
func NewAPIHandler(handler http.Handler, oauthTokenStorage OAuthTokenStorage) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		stok, err := verifySessionToken(req, oauthTokenStorage)

		if err != nil {
			// Erase the session cookie in case of error.
			http.SetCookie(w, &http.Cookie{Name: sessionTokenCookieName, MaxAge: -1})
			w.WriteHeader(http.StatusForbidden)
			fmt.Fprintf(w, "Error: %s.", err)
			return
		}

		// Make sure to refresh the cookie.
		http.SetCookie(w, stok.AsCookie())

		req = req.WithContext(withSessionToken(req.Context(), stok))

		handler.ServeHTTP(w, req)
	})
}

// NewAPIMiddleware instantiates a new API middleware.
func NewAPIMiddleware(storage OAuthTokenStorage) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return NewAPIHandler(handler, storage)
	}
}
