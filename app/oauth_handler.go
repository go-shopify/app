package app

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/go-shopify/shopify"
)

// oauthHandlerImpl represents a Shopify handler.
type oauthHandlerImpl struct {
	Config
	storage      OAuthTokenStorage
	handler      http.Handler
	errorHandler ErrorHandler
}

// NewOAuthHandler instantiates a new Shopify embedded app handler, from the
// specified configuration.
//
// A typical usage of the handler is to serve the `index.html` page of a
// Shopify embedded app.
//
// Upon a successful request, the handler stores or refreshes authentication
// information on the client side, in the form of a cookie.
func NewOAuthHandler(handler http.Handler, storage OAuthTokenStorage, config *Config, errorHandler ErrorHandler) http.Handler {
	if storage == nil {
		panic("An OAuth token storage is required.")
	}

	if config == nil {
		panic("A configuration is required.")
	}

	h := oauthHandlerImpl{
		Config:       *config,
		storage:      storage,
		handler:      handler,
		errorHandler: errorHandler,
	}

	return newHMACHandler(h, h.APISecret)
}

func (h oauthHandlerImpl) handleError(w http.ResponseWriter, req *http.Request, err error) {
	if h.errorHandler != nil {
		h.errorHandler.ServeHTTPError(w, req, err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "Internal server error: you may contact the application adminstrator.\n")
}

func (h oauthHandlerImpl) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	shop := shopify.Shop(req.URL.Query().Get("shop"))

	if shop == "" {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Missing `shop` parameter.")
		return
	}

	state := req.URL.Query().Get("state")

	// If we have a state, assume we are being called back after an install/update.
	if state != "" {
		h.handleInstallationCallback(w, req, shop, state)
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
		h.redirectToInstall(w, req, shop)
		return
	}

	stok := &sessionToken{
		Shop:       shop,
		OAuthToken: *oauthToken,
	}
	http.SetCookie(w, stok.AsCookie())

	req = req.WithContext(withSessionToken(req.Context(), stok))

	h.handler.ServeHTTP(w, req)
}

func (h oauthHandlerImpl) redirectToInstall(w http.ResponseWriter, req *http.Request, shop shopify.Shop) {
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

func (h oauthHandlerImpl) handleInstallationCallback(w http.ResponseWriter, req *http.Request, shop shopify.Shop, state string) {
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

	req = req.WithContext(shopify.WithShop(req.Context(), shop))
	oauthToken, err := shopify.DefaultAdminClient.GetOAuthToken(req.Context(), h.APIKey, h.APISecret, code)

	if err != nil {
		h.handleError(w, req, fmt.Errorf("get OAuth token from Shopify for `%s`: %s", shop, err))
		return
	}

	if err = h.storage.UpdateOAuthToken(req.Context(), shop, *oauthToken); err != nil {
		h.handleError(w, req, fmt.Errorf("updating OAuth token for `%s`: %s", shop, err))
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
