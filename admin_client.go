package shopify

import (
	"net/http"
)

// AdminClient represents a Shopify client, that can interact with the Shopify REST Admin API.
type AdminClient struct {
	// Shop is the shop associated to the admin client.
	Shop Shop

	// HTTPClient is the HTTP client to use for requests.
	//
	// If none is specified, http.DefaultClient is used.
	HTTPClient *http.Client
}

// NewAdminClient instantiates a new admin client for the specified shop.
func NewAdminClient(shop Shop) *AdminClient {
	return &AdminClient{
		Shop: shop,
	}
}

func (c *AdminClient) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	return http.DefaultClient
}
