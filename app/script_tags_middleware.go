package app

import (
	"fmt"
	"net/http"

	"github.com/go-shopify/shopify"
)

// NewScriptTagsMiddleware instantiates a middleware that ensures that a list
// of script tags are registered.
//
// It must be chained with an APIHandler or OAuthHandler as it requires the
// request context to contains the Shopify credentials (shop and access
// tokens).
func NewScriptTagsMiddleware(scriptTags ...shopify.ScriptTag) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			if len(scriptTags) > 0 {
				if _, ok := shopify.GetShop(req.Context()); ok {
					if _, ok := shopify.GetAccessToken(req.Context()); ok {
						for _, scriptTag := range scriptTags {
							if _, err := shopify.DefaultAdminClient.EnsureScriptTag(req.Context(), scriptTag); err != nil {
								w.WriteHeader(http.StatusInternalServerError)
								fmt.Fprintf(w, "The server failed to register script tags. Please contact your administrator.")
								return
							}
						}
					}
				}
			}

			handler.ServeHTTP(w, req)
		})
	}
}
