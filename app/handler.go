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
	handler http.Handler
	storage OAuthTokenStorage
}

// NewHandler instantiates a new Shopify Handler, from the specified
// configuration.
//
// That handler handles OAuth access neogitation and injects shop, locale and
// timestamp information into the request context.
func NewHandler(handler http.Handler, storage OAuthTokenStorage, config *Config) http.Handler {
	if storage == nil {
		panic("An OAuth token storage is required.")
	}

	if config == nil {
		panic("A configuration is required.")
	}

	h := handlerImpl{
		Router:  mux.NewRouter(),
		Config:  *config,
		handler: handler,
		storage: storage,
	}

	h.Router.PathPrefix("/").HandlerFunc(h.delegateOrInstall)

	return NewHMACHandler(h, h.APISecret)
}

func (h handlerImpl) delegateOrInstall(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	state := req.URL.Query().Get("state")

	// If we have a state, assume we are being called back after an install/update.
	if state != "" {
		h.completeInstall(w, req, shop, state)
		return
	}

	oauthToken, err := h.storage.GetOAuthToken(req.Context(), shop)

	if err != nil {
		if h.OnError != nil {
			h.OnError(req.Context(), fmt.Errorf("failed to load OAuth token for `%s`: %s", shop, err))
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	if oauthToken == nil {
		h.redirectToInstall(w, req, shop)
		return
	}

	h.handler.ServeHTTP(w, req)
}

func (h handlerImpl) redirectToInstall(w http.ResponseWriter, req *http.Request, shop shopify.Shop) {
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
	q.Set("scope", h.Scope.String())
	q.Set("state", state)
	q.Set("redirect_uri", h.PublicURL.String())
	oauthURL.RawQuery = q.Encode()

	// Set a cookie to ensure the auth callback is really called by the right
	// entity.
	http.SetCookie(w, &http.Cookie{Name: "state", Value: state})

	// Redirect the browser to the OAuth autorization page.
	clientRedirect(w, req, oauthURL.String())
}

func (h handlerImpl) completeInstall(w http.ResponseWriter, req *http.Request, shop shopify.Shop, state string) {
	stateCookie, err := req.Cookie("state")

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `state` cookie.")
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
	oauthToken, err := adminClient.GetOAuthToken(req.Context(), h.APIKey, h.APISecret, code)

	if err != nil {
		if h.OnError != nil {
			h.OnError(req.Context(), fmt.Errorf("get OAuth token from Shopify for `%s`: %s", shop, err))
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	if err = h.storage.UpdateOAuthToken(req.Context(), shop, *oauthToken); err != nil {
		if h.OnError != nil {
			h.OnError(req.Context(), fmt.Errorf("updating OAuth token for `%s`: %s", shop, err))
		}

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Unexpected error. Please contact the App's administrator.")
		return
	}

	// Remove the state cookie.
	http.SetCookie(w, &http.Cookie{Name: "state", Expires: time.Unix(0, 0)})

	// Redirect the browser to the main page.
	//
	// Make sure parameters are correct or we will redirect to an error page.
	query := url.Values{}
	query.Set("shop", string(shop))
	injectHMAC(query, h.APISecret)

	redirectURL := &url.URL{
		Scheme:   h.PublicURL.Scheme,
		Host:     h.PublicURL.Host,
		Path:     h.PublicURL.Path,
		RawQuery: query.Encode(),
	}

	clientRedirect(w, req, redirectURL.String())
}

func clientRedirect(w http.ResponseWriter, req *http.Request, url string) {
	html := fmt.Sprintf(`
<html>
	<head>
		<script>
			if (window.self === window.top) {
				window.location.href = '%s';
			} else {
				window.top.location.href = '%s';
			}
		</script>
	</head>
</html>
`, url, url)

	w.Header().Set("Content-Type", "text/html")
	fmt.Fprintf(w, "%s", html)
}
