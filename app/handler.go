package app

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-shopify/shopify"
	"github.com/gorilla/mux"
)

// handlerImpl represents a Shopify handler.
type handlerImpl struct {
	*mux.Router
	Config
}

const pathAuthCallback = "/auth/callback"

// NewHandler instantiates a new Shopify Handler, from the specified
// configuration.
func NewHandler(config *Config) http.Handler {
	if config == nil {
		panic("A configuration is required.")
	}

	handler := handlerImpl{
		Router: mux.NewRouter(),
		Config: *config,
	}

	handler.Router.Path("/").Methods(http.MethodGet).Handler(NewHMACHandler(http.HandlerFunc(handler.install), handler.APISecret))
	handler.Router.Path(pathAuthCallback).Methods(http.MethodGet).Handler(NewHMACHandler(http.HandlerFunc(handler.authCallback), handler.APISecret))

	if config.DefaultHandler != nil {
		handler.Router.Handle("", handler.DefaultHandler)
	}

	return handler
}

func (h handlerImpl) install(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	state, err := generateRandomState()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	oauthURL := &url.URL{
		Scheme: "https",
		Host:   string(shop),
		Path:   "/admin/oauth/authorize",
	}

	q := oauthURL.Query()
	q.Set("client_id", string(h.APIKey))
	q.Set("scope", h.Scopes.String())
	q.Set("state", state)
	q.Set("redirect_uri", h.PublicURL.ResolveReference(&url.URL{Path: pathAuthCallback}).String())
	oauthURL.RawQuery = q.Encode()

	// Set a cookie to ensure the auth callback is really called by the right
	// entity.
	http.SetCookie(w, &http.Cookie{Name: "state", Value: state})

	// Redirect the browser to the OAuth autorization page.
	http.Redirect(w, req, oauthURL.String(), http.StatusFound)
}

func (h handlerImpl) authCallback(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	stateCookie, err := req.Cookie("state")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `state` cookie.")
		return
	}

	state := req.URL.Query().Get("state")

	if state == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `state` parameter.")
		return
	}

	if stateCookie.Value != state {
		w.WriteHeader(http.StatusForbidden)
		fmt.Fprintf(w, "A different `state` value was expected.")
		return
	}

	// TODO: Implement.
}
