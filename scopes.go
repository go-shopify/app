package shopify

// Scope represents an OAuth scope, as defined at
// https://help.shopify.com/en/api/getting-started/authentication/oauth/scopes.
type Scope string

const (
	// Authenticated access scopes.

	// ScopeReadContent represents read access to Article, Blog, Comment, Page, and Redirect.
	ScopeReadContent = Scope("read_content")
	// ScopeWriteContent represents write access to Article, Blog, Comment, Page, and Redirect.
	ScopeWriteContent = Scope("write_content")

	// ScopeReadThemes represents read access to Asset and Theme.
	ScopeReadThemes = Scope("read_themes")
	// ScopeWriteThemes represents write access to Asset and Theme.
	ScopeWriteThemes = Scope("write_themes")

	// ScopeReadProducts represents read access to Product, Product Variant, Product Image, Collect, Custom Collection, and Smart Collection.
	ScopeReadProducts = Scope("read_products")
	// ScopeWriteProducts represents write access to Product, Product Variant, Product Image, Collect, Custom Collection, and Smart Collection.
	ScopeWriteProducts = Scope("write_products")

	// ScopeReadProductListings represents read access to Product Listing, and Collection Listing.
	ScopeReadProductListings = Scope("read_product_listings")

	// ScopeReadCustomers represents read access to Customer and Saved Search.
	ScopeReadCustomers = Scope("read_customers")
	// ScopeWriteCustomers represents write access to Customer and Saved Search.
	ScopeWriteCustomers = Scope("write_customers")

	// ScopeReadOrders represents read access to Order, Transaction and Fulfillment.
	ScopeReadOrders = Scope("read_orders")
	// ScopeWriteOrders represents write access to Order, Transaction and Fulfillment.
	ScopeWriteOrders = Scope("write_orders")

	// ScopeReadAllOrders represents read grants access to all orders rather than the default window of 60 days worth of orders. This OAuth scope is used in conjunction with read_orders, or write_orders. You need to request this scope from your Partner Dashboard before adding it to your app.
	ScopeReadAllOrders = Scope("read_all_orders")

	// ScopeReadDraftOrders represents read access to Draft Order.
	ScopeReadDraftOrders = Scope("read_draft_orders")
	// ScopeWriteDraftOrders represents write access to Draft Order.
	ScopeWriteDraftOrders = Scope("write_draft_orders")

	// ScopeReadInventory represents read access to Inventory Level and Inventory Item.
	ScopeReadInventory = Scope("read_inventory")
	// ScopeWriteInventory represents write access to Inventory Level and Inventory Item.
	ScopeWriteInventory = Scope("write_inventory")

	// ScopeReadLocations represents read access to Location.
	ScopeReadLocations = Scope("read_locations")

	// ScopeReadScriptTags represents read access to Script Tag.
	ScopeReadScriptTags = Scope("read_script_tags")

	// ScopeWriteScriptTags represents write access to Script Tag.
	ScopeWriteScriptTags = Scope("write_script_tags")

	// ScopeReadFulfillments represents read access to Fulfillment Service.
	ScopeReadFulfillments = Scope("read_fulfillments")
	// ScopeWriteFulfillments represents write access to Fulfillment Service.
	ScopeWriteFulfillments = Scope("write_fulfillments")

	// ScopeReadShipping represents read access to Carrier Service, Country and Province.
	ScopeReadShipping = Scope("read_shipping")
	// ScopeWriteShipping represents write access to Carrier Service, Country and Province.
	ScopeWriteShipping = Scope("write_shipping")

	// ScopeReadAnalytics represents read access to Analytics API.
	ScopeReadAnalytics = Scope("read_analytics")

	// ScopeReadUsers represents read access to User (SHOPIFY PLUS).
	ScopeReadUsers = Scope("read_users")
	// ScopeWriteUsers represents write access to User (SHOPIFY PLUS).
	ScopeWriteUsers = Scope("write_users")

	// ScopeReadCheckouts represents read access to Checkouts.
	ScopeReadCheckouts = Scope("read_checkouts")
	// ScopeWriteCheckouts represents write access to Checkouts.
	ScopeWriteCheckouts = Scope("write_checkouts")

	// ScopeReadReports represents read access to Reports.
	ScopeReadReports = Scope("read_reports")
	// ScopeWriteReports represents write access to Reports.
	ScopeWriteReports = Scope("write_reports")

	// ScopeReadPriceRules represents read access to Price Rules.
	ScopeReadPriceRules = Scope("read_price_rules")
	// ScopeWritePriceRules represents write access to Price Rules.
	ScopeWritePriceRules = Scope("write_price_rules")

	// ScopeReadMarketingEvents represents read access to Marketing Event.
	ScopeReadMarketingEvents = Scope("read_marketing_events")
	// ScopeWriteMarketingEvents represents write access to Marketing Event.
	ScopeWriteMarketingEvents = Scope("write_marketing_events")

	// ScopeReadResourceFeedbacks represents read access to ResourceFeedback.
	ScopeReadResourceFeedbacks = Scope("read_resource_feedbacks")
	// ScopeWriteResourceFeedbacks represents write access to ResourceFeedback.
	ScopeWriteResourceFeedbacks = Scope("write_resource_feedbacks")

	// ScopeReadShopifyPaymentsPayouts represents read access to Shopify Payments Payouts, and Transactions.
	ScopeReadShopifyPaymentsPayouts = Scope("read_shopify_payments_payouts")

	// Unauthenticated access scopes.

	// ScopeUnauthenticatedReadProductListings represents read unauthenticated access to read the Product and Collection objects.
	ScopeUnauthenticatedReadProductListings = Scope("unauthenticated_read_product_listings")

	// ScopeUnauthenticatedWriteCheckouts represents write unauthenticated access to the Checkout object.
	ScopeUnauthenticatedWriteCheckouts = Scope("unauthenticated_write_checkouts")

	// ScopeUnauthenticatedWriteCustomers represents write unauthenticated access to the Customer object.
	ScopeUnauthenticatedWriteCustomers = Scope("unauthenticated_write_customers")

	// ScopeUnauthenticatedReadCustomerTags represents read unauthenticated access to read the tags field on the Customer object.
	ScopeUnauthenticatedReadCustomerTags = Scope("unauthenticated_read_customer_tags")

	// ScopeUnauthenticatedReadContent represents read unauthenticated access to read storefront content, such as Article, Blog, and Comment objects.
	ScopeUnauthenticatedReadContent = Scope("unauthenticated_read_content")
)
