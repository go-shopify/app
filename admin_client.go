package shopify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
)

// AdminClient represents a Shopify client, that can interact with the Shopify REST Admin API.
type AdminClient struct {
	// Shop is the shop associated to the admin client.
	Shop Shop

	// AccessToken is the access token to use for authentication.
	AccessToken AccessToken

	// HTTPClient is the HTTP client to use for requests.
	//
	// If none is specified, http.DefaultClient is used.
	HTTPClient *http.Client

	shopURL *url.URL
}

const headerXShopifyAccessToken = "X-Shopify-Access-Token"

// NewAdminClient instantiates a new admin client for the specified shop.
func NewAdminClient(shop Shop, accessToken AccessToken) *AdminClient {
	return &AdminClient{
		Shop:        shop,
		AccessToken: accessToken,
		HTTPClient:  newHTTPClient(),
		shopURL: &url.URL{
			Scheme: "https",
			Host:   string(shop),
		},
	}
}

func newHTTPClient() *http.Client {
	client := &http.Client{}

	if Debug {
		client.Transport = &debugTransport{Transport: &http.Transport{}}
	}

	return client
}

func (c *AdminClient) httpClient() *http.Client {
	if c.HTTPClient != nil {
		return c.HTTPClient
	}

	return http.DefaultClient
}

func (c *AdminClient) newURL(path string, values url.Values) *url.URL {
	if values == nil {
		values = url.Values{}
	}

	return c.shopURL.ResolveReference(&url.URL{Path: path, RawQuery: values.Encode()})
}

func (c *AdminClient) newRequest(ctx context.Context, method string, path string, values url.Values, body interface{}) (*http.Request, error) {
	var r io.Reader

	if body != nil {
		data, err := json.Marshal(body)

		if err != nil {
			return nil, fmt.Errorf("failed to JSON-marshal the request body (%#v): %s", body, err)
		}

		r = bytes.NewBuffer(data)
	}

	u := c.newURL(path, values)
	req, err := http.NewRequest(method, u.String(), r)

	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %s", err)
	}

	req = req.WithContext(ctx)

	// If we have a body, assume it will be JSON.
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.AccessToken != "" {
		req.Header.Add(headerXShopifyAccessToken, string(c.AccessToken))
	}

	return req, nil
}

// GetOAuthToken recovers a permanent access token for the associated shop,
// using the specified code.
func (c *AdminClient) GetOAuthToken(ctx context.Context, apiKey APIKey, apiSecret APISecret, code string) (*OAuthToken, error) {
	body := struct {
		ClientID     APIKey    `json:"client_id"`
		ClientSecret APISecret `json:"client_secret"`
		Code         string    `json:"code"`
	}{
		ClientID:     apiKey,
		ClientSecret: apiSecret,
		Code:         code,
	}

	req, err := c.newRequest(ctx, http.MethodPost, "/admin/oauth/access_token", nil, body)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient().Do(req)

	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	defer flushAndCloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
	}

	mediaType, _, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

	if mediaType != "application/json" {
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("unexpected content-type `%s` (body follows):\n%s", mediaType, string(body))
	}

	result := &OAuthToken{}

	if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("unable to parse OAuth token: %s", err)
	}

	return result, nil
}

func flushAndCloseBody(r io.ReadCloser) {
	if r != nil {
		io.Copy(ioutil.Discard, r)
		r.Close()
	}
}
