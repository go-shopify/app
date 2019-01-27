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
	"strconv"
	"strings"
	"sync"
	"time"
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

// Pagination represents pagination options.
type Pagination struct {
	Limit   int
	Page    int
	SinceID int
}

func (o *Pagination) injectInto(values url.Values) {
	if o == nil {
		return
	}

	if o.Limit != 0 {
		values.Set("limit", strconv.Itoa(o.Limit))
	}

	if o.Page != 0 {
		values.Set("page", strconv.Itoa(o.Page))
	}

	if o.SinceID != 0 {
		values.Set("since_id", strconv.Itoa(o.SinceID))
	}
}

const (
	// DefaultLimit is the default limit, as specified by Shopify.
	DefaultLimit = 50

	// MaxLimit is the maximum allowed limit, as specified by Shopify.
	MaxLimit = 250
)

// SelectedFields represents a list of fields to fetch.
type SelectedFields []string

func (f SelectedFields) injectInto(values url.Values) {
	if len(f) == 0 {
		return
	}

	values.Set("fields", strings.Join(f, ","))
}

// ScriptTagEvent represents a script tag event.
type ScriptTagEvent string

const (
	// ScriptTagEventOnLoad is the only possible value.
	ScriptTagEventOnLoad ScriptTagEvent = "onload"
)

// ScriptTagDisplayScope represents a script tag display scope.
type ScriptTagDisplayScope string

const (
	// ScriptTagDisplayScopeOnlineStore indicates that a script tag must be
	// included only on the web storefront.
	ScriptTagDisplayScopeOnlineStore ScriptTagDisplayScope = "online_store"
	// ScriptTagDisplayScopeOrderStatus indicates that a script tag must be
	// included only on the order status page.
	ScriptTagDisplayScopeOrderStatus ScriptTagDisplayScope = "order_status"
	// ScriptTagDisplayScopeAll indicates that a script tag must be
	// included on all pages.
	ScriptTagDisplayScopeAll ScriptTagDisplayScope = "all"
)

// ScriptTagID is an ID of a script tag.
type ScriptTagID int

// ScriptTag represents a script tag.
type ScriptTag struct {
	CreatedAt    time.Time             `json:"created_at,omitempty"`
	Event        ScriptTagEvent        `json:"event"`
	ID           ScriptTagID           `json:"id"`
	Src          string                `json:"src"`
	DisplayScope ScriptTagDisplayScope `json:"display_scope,omitempty"`
	UpdatedAt    time.Time             `json:"updated_at,omitempty"`
}

// EnsureScriptTag makes sure that a specified script tag is registered in the shop.
//
// If the scriptTag has an ID, an optimistic GET is attempted first. If the GET
// succeeds and the script tag is identical, the function exits immediately. No
// duplicates are removed in that case.
//
// Otherwise, all script tags are fetched and compared to the specified one.
// The first script tag that matches exactly is kept (and returned). Any
// additional duplicate script tag is deleted. If no exact match is found, a
// new script tag is created.
func (c *AdminClient) EnsureScriptTag(ctx context.Context, scriptTag ScriptTag) (*ScriptTag, error) {
	if scriptTag.ID != 0 {
		result, err := c.GetScriptTag(ctx, scriptTag.ID, nil)

		if err != nil {
			return nil, fmt.Errorf("failed to lookup existing script tag with ID `%d`: %s", scriptTag.ID, err)
		}

		if result != nil {
			return result, nil
		}
	}

	scriptTags, err := c.GetAllScriptTags(ctx, nil)

	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	defer wg.Wait()

	for _, s := range scriptTags {
		if s.Src == scriptTag.Src {
			if s.DisplayScope == scriptTag.DisplayScope && scriptTag.ID == 0 {
				scriptTag = s
				continue
			}

			wg.Add(1)

			go func(id ScriptTagID) {
				defer wg.Done()
				c.DeleteScriptTag(ctx, id)
			}(s.ID)
		}
	}

	// If we didn't find any existing matching script tag, create one and return it.
	if scriptTag.ID == 0 {
		return c.CreateOrUpdateScriptTag(ctx, scriptTag)
	}

	return &scriptTag, nil
}

// GetAllScriptTags retrieves a list of all script tags.
func (c *AdminClient) GetAllScriptTags(ctx context.Context, fields SelectedFields) ([]ScriptTag, error) {
	count, err := c.GetScriptTagsCount(ctx)

	if err != nil {
		return nil, fmt.Errorf("failed to count script tags: %s", err)
	}

	if count == 0 {
		return nil, nil
	}

	pagination := &Pagination{
		Limit: MaxLimit,
	}

	var result []ScriptTag

	pageCount := ((count - 1) / pagination.Limit) + 1

	for page := 1; page <= pageCount; page++ {
		pagination.Page = page

		scriptTags, err := c.GetScriptTags(ctx, pagination, fields)

		if err != nil {
			return nil, fmt.Errorf("fetching page %d/%d: %s", page, pageCount, err)
		}

		result = append(result, scriptTags...)
	}

	return result, nil
}

// GetScriptTags retrieves a list of script tags.
//
// To fetch the complete list, use GetAllScriptTags.
func (c *AdminClient) GetScriptTags(ctx context.Context, pagination *Pagination, fields SelectedFields) ([]ScriptTag, error) {
	values := url.Values{}
	pagination.injectInto(values)
	fields.injectInto(values)

	req, err := c.newRequest(ctx, http.MethodGet, "/admin/script_tags.json", values, nil)

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

	result := &struct {
		ScriptTags []ScriptTag `json:"script_tags"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return nil, fmt.Errorf("unable to parse response: %s", err)
	}

	return result.ScriptTags, nil
}

// GetScriptTagsCount retrieves the count of all script tags.
func (c *AdminClient) GetScriptTagsCount(ctx context.Context) (int, error) {
	req, err := c.newRequest(ctx, http.MethodGet, "/admin/script_tags/count.json", nil, nil)

	if err != nil {
		return 0, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient().Do(req)

	if err != nil {
		return 0, fmt.Errorf("request failed: %s", err)
	}

	defer flushAndCloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return 0, fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
	}

	result := &struct {
		Count int `json:"count"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(result); err != nil {
		return 0, fmt.Errorf("unable to parse response: %s", err)
	}

	return result.Count, nil
}

// GetScriptTag fetches a script tag by ID.
//
// If no such script tag exists, a nil script tag and no error is returned.
func (c *AdminClient) GetScriptTag(ctx context.Context, id ScriptTagID, fields SelectedFields) (*ScriptTag, error) {
	values := url.Values{}
	fields.injectInto(values)

	req, err := c.newRequest(ctx, http.MethodGet, fmt.Sprintf("/admin/script_tags/%d.json", id), values, nil)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient().Do(req)

	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	defer flushAndCloseBody(resp.Body)

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, nil
	default:
		body, _ := ioutil.ReadAll(resp.Body)

		return nil, fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
	}

	body := &struct {
		ScriptTag ScriptTag `json:"script_tag"`
	}{}

	if err = json.NewDecoder(resp.Body).Decode(body); err != nil {
		return nil, fmt.Errorf("unable to parse response: %s", err)
	}

	return &body.ScriptTag, nil
}

// CreateOrUpdateScriptTag creates or updates a script tag.
//
// If the specified script tag has an ID, an update is attempted.
func (c *AdminClient) CreateOrUpdateScriptTag(ctx context.Context, scriptTag ScriptTag) (*ScriptTag, error) {
	if scriptTag.Event == "" {
		scriptTag.Event = ScriptTagEventOnLoad
	}

	body := &struct {
		ScriptTag ScriptTag `json:"script_tag"`
	}{
		ScriptTag: scriptTag,
	}

	var req *http.Request
	var err error

	if scriptTag.ID == 0 {
		req, err = c.newRequest(ctx, http.MethodPost, "/admin/script_tags.json", nil, body)
	} else {
		req, err = c.newRequest(ctx, http.MethodPut, fmt.Sprintf("/admin/script_tags/%d.json", scriptTag.ID), nil, body)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient().Do(req)

	if err != nil {
		return nil, fmt.Errorf("request failed: %s", err)
	}

	defer flushAndCloseBody(resp.Body)

	if scriptTag.ID == 0 {
		if resp.StatusCode != http.StatusCreated {
			body, _ := ioutil.ReadAll(resp.Body)

			return nil, fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
		}
	} else {
		if resp.StatusCode != http.StatusOK {
			body, _ := ioutil.ReadAll(resp.Body)

			return nil, fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
		}
	}

	if err = json.NewDecoder(resp.Body).Decode(body); err != nil {
		return nil, fmt.Errorf("unable to parse response: %s", err)
	}

	return &body.ScriptTag, nil
}

// DeleteScriptTag deletes a script tag.
func (c *AdminClient) DeleteScriptTag(ctx context.Context, id ScriptTagID) error {
	req, err := c.newRequest(ctx, http.MethodDelete, fmt.Sprintf("/admin/script_tags/%d.json", id), nil, nil)

	if err != nil {
		return fmt.Errorf("failed to create request: %s", err)
	}

	resp, err := c.httpClient().Do(req)

	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}

	defer flushAndCloseBody(resp.Body)

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)

		return fmt.Errorf("unexpected return status code of %d (body follows):\n%s", resp.StatusCode, string(body))
	}

	return nil
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
