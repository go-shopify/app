package shopify

// OAuthToken represents an OAuth token as received from a shop.
type OAuthToken struct {
	AccessToken AccessToken `json:"access_token"`
	Scope       Scope       `json:"scope"`
}
