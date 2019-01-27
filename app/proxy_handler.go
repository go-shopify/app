package app

import (
	"fmt"
	"net/http"

	"github.com/go-shopify/shopify"
)

type proxyHandlerImpl struct {
	Config
	storage      OAuthTokenStorage
	handler      http.Handler
	errorHandler ErrorHandler
}

// NewProxyHandler instantiates a new Shopify proxy handler, from the specified
// configuration.
//
// A typical usage of the handler is to serve pages, scripts or APIs through a
// Shopify App proxy, usually from the storefront.
func NewProxyHandler(handler http.Handler, storage OAuthTokenStorage, config *Config, errorHandler ErrorHandler) http.Handler {
	if storage == nil {
		panic("An OAuth token storage is required.")
	}

	if config == nil {
		panic("A configuration is required.")
	}

	h := proxyHandlerImpl{
		Config:       *config,
		storage:      storage,
		handler:      handler,
		errorHandler: errorHandler,
	}

	return newSignatureHandler(h, h.APISecret)
}

func (h proxyHandlerImpl) handleError(w http.ResponseWriter, req *http.Request, err error) {
	if h.errorHandler != nil {
		h.errorHandler.ServeHTTPError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal server error: you may contact the application adminstrator.\n")
}

func (h proxyHandlerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	// Load any existing OAuth token for that shop.
	oauthToken, err := h.storage.GetOAuthToken(req.Context(), shop)

	if err != nil {
		h.handleError(w, req, fmt.Errorf("failed to load OAuth token for `%s`: %s", shop, err))
		return
	}

	// If we don't have a token yet for that shop, redirect for the OAuth page.
	if oauthToken == nil {
		h.handleError(w, req, fmt.Errorf("no OAuth token for `%s`", shop))
		return
	}

	stok := &sessionToken{
		Shop:       shop,
		OAuthToken: *oauthToken,
	}

	req = req.WithContext(withSessionToken(req.Context(), stok))

	h.handler.ServeHTTP(w, req)
}

// NewProxyMiddleware instantiates a new proxy middleware.
func NewProxyMiddleware(storage OAuthTokenStorage, config *Config, errorHandler ErrorHandler) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return NewProxyHandler(handler, storage, config, errorHandler)
	}
}
