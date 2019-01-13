package app

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

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
//
// That handler handles OAuth access neogitation and injects shop, locale and
// timestamp information into the request context.
func NewHandler(config *Config) http.Handler {
	if config == nil {
		panic("A configuration is required.")
	}

	handler := handlerImpl{
		Router: mux.NewRouter(),
		Config: *config,
	}

	handler.Router.Path(pathAuthCallback).Methods(http.MethodGet).HandlerFunc(handler.authCallback)
	handler.Router.PathPrefix("/").HandlerFunc(handler.delegateOrInstall)

	return NewHMACHandler(handler, handler.APISecret)
}

func (h handlerImpl) delegateOrInstall(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	req = req.WithContext(WithShop(req.Context(), shop))

	accessToken, err := h.OnAccessTokenRequested(req.Context(), shop)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	if accessToken == "" {
		h.install(w, req, shop)
		return
	}

	if h.Handler != nil {
		h.Handler.ServeHTTP(w, req)
	}
}

func (h handlerImpl) install(w http.ResponseWriter, req *http.Request, shop shopify.Shop) {
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

	code := req.URL.Query().Get("code")

	if code == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `code` parameter.")
		return
	}

	adminClient := shopify.NewAdminClient(shop, "")
	accessToken, err := adminClient.GetOAuthAccessToken(req.Context(), h.APIKey, h.APISecret, code)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	if h.OnAccessTokenUpdated != nil {
		if err = h.OnAccessTokenUpdated(req.Context(), shop, accessToken); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
			return
		}
	}

	// Remove the state cookie.
	http.SetCookie(w, &http.Cookie{Name: "state", Expires: time.Unix(0, 0)})

	// Redirect the browser to the main page.
	http.Redirect(w, req, h.PublicURL.String(), http.StatusFound)
}
