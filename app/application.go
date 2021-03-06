package app

import (
	"net/http"

	"github.com/go-shopify/shopify"
)

// Application represents a Shopify embedded application.
type Application struct {
	// Config contains the application configuration.
	Config *Config

	// OAuthTokenStorage contains the OAuth token storage provider.
	OAuthTokenStorage OAuthTokenStorage

	// ErrorHandler, if specified, is a handler to call when fatal errors occurs.
	ErrorHandler ErrorHandler
}

// NewOAuthHandler instantiates a new Shopify embedded app handler.
//
// A typical usage of the handler is to serve the `index.html` page of a
// Shopify embedded app.
//
// Upon a successful request, the handler stores or refreshes authentication
// information on the client side, in the form of a cookie.
func (a *Application) NewOAuthHandler(handler http.Handler) http.Handler {
	return NewOAuthHandler(handler, a.OAuthTokenStorage, a.Config, a.ErrorHandler)
}

// NewOAuthMiddleware instantiates a new Shopify embedded app middleware.
//
// A typical usage of the handler is to serve the `index.html` page of a
// Shopify embedded app.
//
// Upon a successful request, the handler stores or refreshes authentication
// information on the client side, in the form of a cookie.
func (a *Application) NewOAuthMiddleware() func(http.Handler) http.Handler {
	return NewOAuthMiddleware(a.OAuthTokenStorage, a.Config, a.ErrorHandler)
}

// NewScriptTagsHandler instantiates a new script tags handler.
//
// A typical usage of the handler is to serve the `index.html` page of a
// Shopify embedded app.
//
// Upon a successful request, the handler stores or refreshes authentication
// information on the client side, in the form of a cookie.
func (a *Application) NewScriptTagsHandler(handler http.Handler, scriptTags ...shopify.ScriptTag) http.Handler {
	return NewScriptTagsHandler(handler, scriptTags...)
}

// NewScriptTagsMiddleware instantiates a new script tags middleware.
//
// A typical usage of the handler is to serve the `index.html` page of a
// Shopify embedded app.
//
// Upon a successful request, the handler stores or refreshes authentication
// information on the client side, in the form of a cookie.
func (a *Application) NewScriptTagsMiddleware(scriptTags ...shopify.ScriptTag) func(http.Handler) http.Handler {
	return NewScriptTagsMiddleware(scriptTags...)
}

// NewProxyHandler instantiates a new Shopify proxy handler.
//
// A typical usage of the handler is to serve pages, scripts or APIs through a
// Shopify App proxy, usually from the storefront.
func (a *Application) NewProxyHandler(handler http.Handler) http.Handler {
	return NewProxyHandler(handler, a.OAuthTokenStorage, a.Config, a.ErrorHandler)
}

// NewProxyMiddleware instantiates a new Shopify proxy handler.
//
// A typical usage of the handler is to serve pages, scripts or APIs through a
// Shopify App proxy, usually from the storefront.
func (a *Application) NewProxyMiddleware() func(http.Handler) http.Handler {
	return NewProxyMiddleware(a.OAuthTokenStorage, a.Config, a.ErrorHandler)
}

// NewAPIHandler instantiates a new API handler.
//
// A typical usage is to wrap custom API rest endpoints with an APIHandler to
// ensure that the calls originates from a Shopify admin page that went through
// a OAuthHandler.
func (a *Application) NewAPIHandler(handler http.Handler) http.Handler {
	return NewAPIHandler(handler, a.OAuthTokenStorage)
}

// NewAPIMiddleware instantiates a new API middleware.
//
// A typical usage is to wrap custom API rest endpoints with an APIHandler to
// ensure that the calls originates from a Shopify admin page that went through
// a OAuthHandler.
func (a *Application) NewAPIMiddleware() func(http.Handler) http.Handler {
	return NewAPIMiddleware(a.OAuthTokenStorage)
}
