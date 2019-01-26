package shopify

import (
	"encoding/json"
	"strings"
)

// Permission represents an OAuth scope, as defined at
// https://help.shopify.com/en/api/getting-started/authentication/oauth/scopes.
type Permission string

// Scope represents a list of permissions.
type Scope []Permission

func (s Scope) String() string {
	strs := make([]string, len(s))

	for i, perm := range s {
		strs[i] = string(perm)
	}

	return strings.Join(strs, ",")
}

// Equal compares two scopes.
func (s Scope) Equal(other Scope) bool {
	if len(s) != len(other) {
		return false
	}

	for i, v := range s {
		if other[i] != v {
			return false
		}
	}

	return true
}

// ParseScope parses a string of comma-separated permissions.
//
// If the string has an incorrect format, an error is returned.
func ParseScope(s string) (result Scope, err error) {
	// Note: We don't currently return any error, but it may happen so we
	// already make the function return an error to prevent future API breaks.

	strs := strings.Split(s, ",")

	result = make(Scope, 0, len(strs))

	for _, perm := range strs {
		perm = strings.TrimSpace(perm)

		if perm != "" {
			result = append(result, Permission(perm))
		}
	}

	return
}

// MarshalJSON implements JSON marshalling.
func (s Scope) MarshalJSON() ([]byte, error) {
	return json.Marshal(s.String())
}

// UnmarshalJSON implements JSON unmarshalling
func (s *Scope) UnmarshalJSON(b []byte) (err error) {
	var str string

	if err = json.Unmarshal(b, &str); err != nil {
		return err
	}

	*s, err = ParseScope(str)

	return
}

const (
	// Authenticated access scopes.

	// PermissionReadContent represents read access to Article, Blog, Comment, Page, and Redirect.
	PermissionReadContent = Permission("read_content")
	// PermissionWriteContent represents write access to Article, Blog, Comment, Page, and Redirect.
	PermissionWriteContent = Permission("write_content")

	// PermissionReadThemes represents read access to Asset and Theme.
	PermissionReadThemes = Permission("read_themes")
	// PermissionWriteThemes represents write access to Asset and Theme.
	PermissionWriteThemes = Permission("write_themes")

	// PermissionReadProducts represents read access to Product, Product Variant, Product Image, Collect, Custom Collection, and Smart Collection.
	PermissionReadProducts = Permission("read_products")
	// PermissionWriteProducts represents write access to Product, Product Variant, Product Image, Collect, Custom Collection, and Smart Collection.
	PermissionWriteProducts = Permission("write_products")

	// PermissionReadProductListings represents read access to Product Listing, and Collection Listing.
	PermissionReadProductListings = Permission("read_product_listings")

	// PermissionReadCustomers represents read access to Customer and Saved Search.
	PermissionReadCustomers = Permission("read_customers")
	// PermissionWriteCustomers represents write access to Customer and Saved Search.
	PermissionWriteCustomers = Permission("write_customers")

	// PermissionReadOrders represents read access to Order, Transaction and Fulfillment.
	PermissionReadOrders = Permission("read_orders")
	// PermissionWriteOrders represents write access to Order, Transaction and Fulfillment.
	PermissionWriteOrders = Permission("write_orders")

	// PermissionReadAllOrders represents read grants access to all orders rather than the default window of 60 days worth of orders. This OAuth scope is used in conjunction with read_orders, or write_orders. You need to request this scope from your Partner Dashboard before adding it to your app.
	PermissionReadAllOrders = Permission("read_all_orders")

	// PermissionReadDraftOrders represents read access to Draft Order.
	PermissionReadDraftOrders = Permission("read_draft_orders")
	// PermissionWriteDraftOrders represents write access to Draft Order.
	PermissionWriteDraftOrders = Permission("write_draft_orders")

	// PermissionReadInventory represents read access to Inventory Level and Inventory Item.
	PermissionReadInventory = Permission("read_inventory")
	// PermissionWriteInventory represents write access to Inventory Level and Inventory Item.
	PermissionWriteInventory = Permission("write_inventory")

	// PermissionReadLocations represents read access to Location.
	PermissionReadLocations = Permission("read_locations")

	// PermissionReadScriptTags represents read access to Script Tag.
	PermissionReadScriptTags = Permission("read_script_tags")

	// PermissionWriteScriptTags represents write access to Script Tag.
	PermissionWriteScriptTags = Permission("write_script_tags")

	// PermissionReadFulfillments represents read access to Fulfillment Service.
	PermissionReadFulfillments = Permission("read_fulfillments")
	// PermissionWriteFulfillments represents write access to Fulfillment Service.
	PermissionWriteFulfillments = Permission("write_fulfillments")

	// PermissionReadShipping represents read access to Carrier Service, Country and Province.
	PermissionReadShipping = Permission("read_shipping")
	// PermissionWriteShipping represents write access to Carrier Service, Country and Province.
	PermissionWriteShipping = Permission("write_shipping")

	// PermissionReadAnalytics represents read access to Analytics API.
	PermissionReadAnalytics = Permission("read_analytics")

	// PermissionReadUsers represents read access to User (SHOPIFY PLUS).
	PermissionReadUsers = Permission("read_users")
	// PermissionWriteUsers represents write access to User (SHOPIFY PLUS).
	PermissionWriteUsers = Permission("write_users")

	// PermissionReadCheckouts represents read access to Checkouts.
	PermissionReadCheckouts = Permission("read_checkouts")
	// PermissionWriteCheckouts represents write access to Checkouts.
	PermissionWriteCheckouts = Permission("write_checkouts")

	// PermissionReadReports represents read access to Reports.
	PermissionReadReports = Permission("read_reports")
	// PermissionWriteReports represents write access to Reports.
	PermissionWriteReports = Permission("write_reports")

	// PermissionReadPriceRules represents read access to Price Rules.
	PermissionReadPriceRules = Permission("read_price_rules")
	// PermissionWritePriceRules represents write access to Price Rules.
	PermissionWritePriceRules = Permission("write_price_rules")

	// PermissionReadMarketingEvents represents read access to Marketing Event.
	PermissionReadMarketingEvents = Permission("read_marketing_events")
	// PermissionWriteMarketingEvents represents write access to Marketing Event.
	PermissionWriteMarketingEvents = Permission("write_marketing_events")

	// PermissionReadResourceFeedbacks represents read access to ResourceFeedback.
	PermissionReadResourceFeedbacks = Permission("read_resource_feedbacks")
	// PermissionWriteResourceFeedbacks represents write access to ResourceFeedback.
	PermissionWriteResourceFeedbacks = Permission("write_resource_feedbacks")

	// PermissionReadShopifyPaymentsPayouts represents read access to Shopify Payments Payouts, and Transactions.
	PermissionReadShopifyPaymentsPayouts = Permission("read_shopify_payments_payouts")

	// Unauthenticated access scopes.

	// PermissionUnauthenticatedReadProductListings represents read unauthenticated access to read the Product and Collection objects.
	PermissionUnauthenticatedReadProductListings = Permission("unauthenticated_read_product_listings")

	// PermissionUnauthenticatedWriteCheckouts represents write unauthenticated access to the Checkout object.
	PermissionUnauthenticatedWriteCheckouts = Permission("unauthenticated_write_checkouts")

	// PermissionUnauthenticatedWriteCustomers represents write unauthenticated access to the Customer object.
	PermissionUnauthenticatedWriteCustomers = Permission("unauthenticated_write_customers")

	// PermissionUnauthenticatedReadCustomerTags represents read unauthenticated access to read the tags field on the Customer object.
	PermissionUnauthenticatedReadCustomerTags = Permission("unauthenticated_read_customer_tags")

	// PermissionUnauthenticatedReadContent represents read unauthenticated access to read storefront content, such as Article, Blog, and Comment objects.
	PermissionUnauthenticatedReadContent = Permission("unauthenticated_read_content")
)
