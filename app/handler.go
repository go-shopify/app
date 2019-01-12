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

	handler.Router.Path("/").Methods(http.MethodGet).HandlerFunc(handler.install)
	handler.Router.Path("/auth/callback").Methods(http.MethodGet).HandlerFunc(handler.install)

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
	q.Set("redirect_uri", "/auth/callback")
	oauthURL.RawQuery = q.Encode()

	// Set a cookie to ensure the auth callback is really called by the right
	// entity.
	http.SetCookie(w, &http.Cookie{Name: "state", Value: state})
	http.Redirect(w, req, oauthURL.String(), http.StatusFound)
}

func (h handlerImpl) authCallback(w http.ResponseWriter, req *http.Request) {
	// TODO: Implement.
}
